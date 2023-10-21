import yaml
from sentence_transformers import SentenceTransformer
from yaml import Loader

config_path = './app/search_vector/config/config.yaml'

def load_website():
    with open(config_path, 'r') as f:
        conf = yaml.load(f, Loader=Loader)

    websites_conf = conf['websites']
    return websites_conf['host'], websites_conf['port']


def load_milvus():
    with open(config_path, 'r') as f:
        conf = yaml.load(f, Loader=Loader)
    milvus_conf = conf['milvus']
    return milvus_conf['host'], milvus_conf['port'], milvus_conf['default_milvus_table_name'], milvus_conf[
        'vector_dimension'], milvus_conf['metric_type']


def load_model():
    with open(config_path, 'r') as f:
        conf = yaml.load(f, Loader=Loader)
    model_conf = conf['model']
    return model_conf['sentence_transformer'], model_conf['network']


MILVUS_HOST, MILVUS_PORT, DEFAULT_MILVUS_TABLE_NAME, VECTOR_DIMENSION, METRIC_TYPE = load_milvus()
WEBSITE_HOST, WEBSITE_PORT = load_website()
TRANSFORMER_MODEL_NAME, NETWORK_MODEL_NAME = load_model()
TRANSFORMER_MODEL = SentenceTransformer(TRANSFORMER_MODEL_NAME)


def init_config_test():
    conf = load_website()
    print(conf)


if __name__ == '__main__':
    init_config_test()
