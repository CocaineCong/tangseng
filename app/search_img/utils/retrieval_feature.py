import os
import pickle

import torch
from PIL import Image
from torchvision import transforms

from cirtorch.networks.imageretrievalnet import init_network, extract_vectors
from lshash.lshash import LSHash

# setting up the visible GPU
os.environ['CUDA_VISIBLE_DEVICES'] = "0"


class ImageProcess:
    def __init__(self, img_dir):
        self.img_dir = img_dir

    def process(self):
        imgs = list()
        for root, dirs, files in os.walk(self.img_dir):
            for file in files:
                img_path = os.path.join(root + os.sep, file)
                try:
                    image = Image.open(img_path)
                    if max(image.size) / min(image.size) < 5:
                        imgs.append(img_path)
                    else:
                        continue
                except:
                    print("image height/width ratio is small")

        return imgs


class AntiFraudFeatureDataset:
    def __init__(self, img_dir, network, feature_path='', index_path=''):
        self.img_dir = img_dir
        self.network = network
        self.feature_path = feature_path
        self.index_path = index_path

    def construct_feature(self, hash_size, input_dim, num_hash_tables):
        multiscale = '[1]'
        print(">> Loading network:\n>>>> '{}'".format(self.network))
        state = torch.load(self.network)
        # parsing net params from meta
        # architecture, pooling, mean, std required
        # the rest has default values, in case that is doesnt exist
        net_params = {
            'architecture': state['meta']['architecture'],
            'pooling': state['meta']['pooling'],
            'local_whitening': state['meta'].get('local_whitening', False),
            'regional': state['meta'].get('regional', False),
            'whitening': state['meta'].get('whitening', False),
            'mean': state['meta']['mean'],
            'std': state['meta']['std'],
            'pretrained': False
        }
        # network initialization
        net = init_network(net_params)
        net.load_state_dict(state['state_dict'])
        # setting up the multi-scale parameters
        ms = list(eval(multiscale))
        print(">>>> Evaluating scales: {}".format(ms))
        # moving network to gpu and eval mode
        if torch.cuda.is_available():
            net.cuda()
        net.eval()

        # set up the transform
        normalize = transforms.Normalize(
            mean=net.meta['mean'],
            std=net.meta['std']
        )
        transform = transforms.Compose([
            transforms.ToTensor(),
            normalize
        ])

        # extract database and query vectors
        print('>> database images...')
        images = ImageProcess(self.img_dir).process()
        vecs, img_paths = extract_vectors(net, images, 1024, transform, ms=ms)
        img_paths = [b for a in img_paths for b in a]
        feature_dict = dict(zip(img_paths, list(vecs.detach().cpu().numpy().T)))
        # index
        lsh = LSHash(hash_size=int(hash_size), input_dim=int(input_dim), num_hashtables=int(num_hash_tables))
        for img_path, vec in feature_dict.items():
            lsh.index(vec.flatten(), extra_data=img_path)

        ## Save Index Model
        with open(self.feature_path, "wb") as f:
            pickle.dump(feature_dict, f)
        with open(self.index_path, "wb") as f:
            pickle.dump(lsh, f)

        print("extract feature is done")
        return feature_dict, lsh
