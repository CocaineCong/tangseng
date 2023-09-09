import io
import json
import os
import pickle
import sys
from base64 import encodebytes

import numpy as np
import torch
import yaml
from PIL import Image
from flask import Flask, request
from torchvision import transforms
from yaml import Loader

from cirtorch.datasets.datahelpers import imresize
from cirtorch.networks.imageretrievalnet import init_network

app = Flask(__name__)


# the entrance of the flask
@app.route("/image", methods=['POST'])
def accInsurance():
    try:
        app.logger.debug("print headers------")
        headers = request.headers
        headers_info = ""
        for k, v in headers.items():
            headers_info += "{}: {}\n".format(k, v)
        app.logger.debug(headers_info)

        app.logger.debug("print forms------")
        forms_info = ""
        for k, v in request.form.items():
            forms_info += "{}: {}\n".format(k, v)
        app.logger.debug(forms_info)

        if 'query' not in request.files:
            return json.dumps({'err': 2, 'msg': 'query image is empty'})

        img_name = request.files['query'].filename
        img_bytes = request.files['query'].read()
        img = request.files['query']

        if img_bytes is None:
            return json.dumps({'err': 3, 'msg': 'img is none'})

        results = retrieval(img)

        data = dict()
        data['query'] = img_name
        data['results'] = results

        return json.dumps({'err': 0, 'msg': 'success', 'data': data})
    except Exception as e:
        app.logger.exception(sys.exc_info())
        return json.dumps({'err': 1, 'msg': e})


# Get the image encoded by Base64
def get_response_image(image_path):
    pil_image = Image.open(image_path, mode='r').convert("RGB")
    byte_arr = io.BytesIO()
    pil_image.save(byte_arr, format='JPEG')
    encoded_img = encodebytes(byte_arr.getvalue()).decode('ascii')
    return encoded_img


# Compute the cosine relativity
def cosine_dist(x, y):
    return 100 * float(np.dot(x, y)) / (np.dot(x, x) * np.dot(y, y)) ** 0.5


# Inference from ResNet-50
def inference(img):
    try:
        # Preprocess the image
        input = Image.open(img).convert("RGB")
        input = imresize(input, 224)
        input = transforms(input).unsqueeze(0)
        if torch.cuda.is_available():
            input = input.cuda()
        # Perform the prediction
        with torch.no_grad():
            vect = net(input)
        return vect
    except Exception as e:
        print(e)


# Use LSH to Query the similar image
def retrieval(img):
    # load model
    query_vect = inference(img)
    query_vect = query_vect.cpu()
    query_vect = list(query_vect.detach().numpy().T[0])

    # LSH Query
    # In order to speed up the reaction, we just get 3 of the most similar image
    response = lsh.query(query_vect, num_results=3, distance_func="cosine")
    try:
        results = []
        for i in range(3):
            similar_path = response[i][0][1]
            # compute the relativity of the query image and the result image
            score = np.rint(cosine_dist(list(query_vect), list(response[i][0][0])))
            result = {"score": score, "image": get_response_image(similar_path)}
            results.append(result)
    except Exception as e:
        results = []
        print(e)

    return results


def init_model(network):
    print(">> Loading network:\n>>>> '{}'".format(network))
    state = torch.load(network)
    # parsing net params from meta
    # architecture, pooling, mean, std required
    # the rest has default values, in case that is doesnt exist
    net_params = {
        'architecture': state['meta']['architecture'],
        'pooling': state['meta']['pooling'],
        'local_whitening': state['meta'].get('local_whitening', True),
        'regional': state['meta'].get('regional', False),
        'whitening': state['meta'].get('whitening', False),
        'mean': state['meta']['mean'],
        'std': state['meta']['std'],
        'pretrained': False
    }
    # network initialization
    net = init_network(net_params)
    net.load_state_dict(state['state_dict'])
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

    with open(os.path.join("./index/", "dataset_index_wukong.pkl"), "rb") as f:
        lsh = pickle.load(f)

    return net, lsh, transform


def init():
    with open('config.yaml', 'r') as f:
        conf = yaml.load(f, Loader=Loader)

    app.logger.info(conf)
    host = conf['websites']['host']
    port = conf['websites']['port']
    network = conf['model']['network']

    net, lsh, transforms = init_model(network)

    return host, port, net, lsh, transforms


if __name__ == "__main__":
    host, port, net, lsh, transforms = init()
    app.run(host=host, port=port, debug=True)
    print("start server {}:{}".format(host, port))
