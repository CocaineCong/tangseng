import yaml
from yaml import Loader
from image_retrieval import init_model

VECTOR_DIMENSION = 3
MILVUS_HOTS = "localhost"
MILVUS_PORT = "19530"
DEFAULT_MILVUS_TABLE_NAME = "milvus_table_name"


def init_config():
    with open('./config.yaml', 'r') as f:
        conf = yaml.load(f, Loader=Loader)

    host = conf['websites']['host']
    port = conf['websites']['port']
    network = conf['model']['network']

    MILVUS_HOTS = conf['milvus']['host']
    MILVUS_PORT = conf['milvus']['port']

    net, lsh, transforms = init_model(network)

    return host, port, net, lsh, transforms

