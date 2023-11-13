import asyncio
import io
import json
import os
import pickle
import sys
from base64 import encodebytes

import numpy as np
import torch
from PIL import Image
from flask import Flask, request
from torchvision import transforms
from app.search_vector.service.search_vector import serve

from app.search_vector.config.config import DEFAULT_MILVUS_TABLE_NAME, NETWORK_MODEL_NAME
from app.search_vector.cirtorch.datasets.datahelpers import imresize
from app.search_vector.cirtorch.networks.imageretrievalnet import init_network
from app.search_vector.milvus.milvus import milvus_client
from app.search_vector.milvus.operators import do_upload, do_search
from app.search_vector.utils.logs import LOGGER

app = Flask(__name__)


@app.route("/test_insert", methods=['GET'])
def test_insert_something():
    ids = do_upload(DEFAULT_MILVUS_TABLE_NAME, 2, "mirror",
                    "mirror something like mirror", milvus_client)
    print(ids)
    return json.dumps({'err': 0, 'msg': 'ok', 'data': 'ok'})


@app.route("/test_search", methods=['POST'])
def test_search_something():
    query = request.form.get('query')
    print(query)
    ids, distance = do_search(DEFAULT_MILVUS_TABLE_NAME, query, 1,
                              milvus_client)
    print(ids)
    print(distance)
    return json.dumps({'err': 0, 'msg': 'ok', 'data': 'ok'})


# the entrance of the flask
@app.route("/image", methods=['POST'])
def accInsurance():
    try:
        LOGGER.debug("print headers------")
        headers = request.headers
        headers_info = ""
        for k, v in headers.items():
            headers_info += "{}: {}\n".format(k, v)
        LOGGER.debug(headers_info)

        LOGGER.debug("print forms------")
        forms_info = ""
        for k, v in request.form.items():
            forms_info += "{}: {}\n".format(k, v)
        LOGGER.debug(forms_info)

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
        LOGGER.exception(sys.exc_info())
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
    return 100 * float(np.dot(x, y)) / (np.dot(x, x) * np.dot(y, y))**0.5


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
            score = np.rint(
                cosine_dist(list(query_vect), list(response[i][0][0])))
            result = {
                "score": score,
                "image": get_response_image(similar_path)
            }
            results.append(result)
    except Exception as e:
        results = []
        print(e)

    return results


def init_model():
    network = NETWORK_MODEL_NAME
    print(">> Loading network:\n>>>> '{}'".format(network))
    state = torch.load(network)
    # parsing net params from meta
    # architecture, pooling, mean, std required
    # the rest has default values, in case that is doesn't exist
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
    normalize = transforms.Normalize(mean=net.meta['mean'],
                                     std=net.meta['std'])
    transform = transforms.Compose([transforms.ToTensor(), normalize])

    with open(
            os.path.join("app/search_vector/index/",
                         "dataset_index_wukong.pkl"), "rb") as f:
        lsh = pickle.load(f)

    return net, lsh, transform


net, lsh, transform = init_model()

if __name__ == "__main__":
    # app.run(host=WEBSITE_HOST, port=WEBSITE_PORT, debug=True)
    # print("start server {}:{}".format(WEBSITE_HOST, WEBSITE_PORT))
    asyncio.run(serve())
