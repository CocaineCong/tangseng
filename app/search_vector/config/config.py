import yaml
from sentence_transformers import SentenceTransformer
from yaml import Loader

config_path = 'config/config.yaml'


def load_website():
    with open(config_path, 'r') as f:
        conf = yaml.load(f, Loader=Loader)
    
    websites_conf = conf['services']
    return websites_conf['search_vector']['addr'][0]


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


def load_etcd():
    with open(config_path, 'r') as f:
        conf = yaml.load(f, Loader=Loader)
    etcd_conf = conf['etcd']
    etcd_list = str.split(etcd_conf['address'],":")
    return etcd_list[0], etcd_list[1]


MILVUS_HOST, MILVUS_PORT, DEFAULT_MILVUS_TABLE_NAME, VECTOR_DIMENSION, METRIC_TYPE = load_milvus()
VECTOR_ADDR = load_website()
TRANSFORMER_MODEL_NAME, NETWORK_MODEL_NAME = load_model()
TRANSFORMER_MODEL = SentenceTransformer(TRANSFORMER_MODEL_NAME)
ETCD_HOST, ETCD_PORT = load_etcd()


def init_config_test():
    conf = load_website()
    print(conf)


if __name__ == '__main__':
    init_config_test()
