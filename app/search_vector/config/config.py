# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""
loading config information
"""
from sentence_transformers import SentenceTransformer
from yaml import Loader, load

CONFIG_PATH = 'config/config.yaml'


def load_website():
    """
    loading the website infomation from config path.
    return the server address for rpc call.
    """

    with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
        conf = load(f, Loader=Loader)
    websites_conf = conf['services']
    return websites_conf['search_vector']['addr'][0]


def load_milvus():
    """
    read config from path and loading some infomation about milvus
    milvus is a vector database. so cool~
    """

    with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
        conf = load(f, Loader=Loader)
    milvus_conf = conf['milvus']
    return milvus_conf['host'], milvus_conf['port'], milvus_conf[
        'default_milvus_table_name'], milvus_conf[
            'vector_dimension'], milvus_conf['metric_type']


def load_model():
    """
    read config from path and loading some infomation about model
    """

    with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
        conf = load(f, Loader=Loader)
    model_conf = conf['model']
    return model_conf['sentence_transformer'], model_conf['network']


def load_etcd():
    """
    read config from path and loading etcd # TODO maybe we will support cluster later :-)
    """

    with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
        conf = load(f, Loader=Loader)
    etcd_conf = conf['etcd']
    etcd_list = str.split(etcd_conf['address'], ":")
    return etcd_list[0], etcd_list[1]


def load_kafka():
    """
    read config from path and loading kafka cluster

    return kafka cluster list 
    
    eg: return ['localhost:9200','localhost:9300','loadlhost:9400']
    """

    with open(CONFIG_PATH, 'r', encoding='utf-8') as f:
        conf = load(f, Loader=Loader)
    kafka_conf = conf['kafka']
    return kafka_conf['address']


MILVUS_HOST, MILVUS_PORT, DEFAULT_MILVUS_TABLE_NAME, VECTOR_DIMENSION, METRIC_TYPE = load_milvus(
)
VECTOR_ADDR = load_website()
TRANSFORMER_MODEL_NAME, NETWORK_MODEL_NAME = load_model()
TRANSFORMER_MODEL = SentenceTransformer(TRANSFORMER_MODEL_NAME)
ETCD_HOST, ETCD_PORT = load_etcd()
KAFKA_CLUSTER = load_kafka()
