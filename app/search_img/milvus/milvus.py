import sys
from pymilvus import (connections, Collection, utility,
                      FieldSchema, DataType, CollectionSchema)
from config.config import VECTOR_DIMENSION, MILVUS_HOST, MILVUS_PORT
from utils.logs import LOGGER


class Milvus:

    def __init__(self, host=MILVUS_HOST, port=MILVUS_PORT):
        try:
            self.collection = None
            connections.connect(host=host, port=port)
        except Exception as e:
            LOGGER.error(f"init milvus {e}")
            sys.exit(1)

    def set_collection(self, collection_name):
        try:
            self.collection = Collection(name=collection_name)
        except Exception as e:
            LOGGER.error(f"failed to set collection to milvus {e}")
            sys.exit(1)

    def has_collection(self, collection_name):
        try:
            return utility.has_collection(collection_name)
        except Exception as e:
            LOGGER.error(f"failed to has collection to milvus {e}")
            sys.exit(1)

    def create_collection(self, collection_name):
        try:
            if not self.has_collection(collection_name):
                doc_id = FieldSchema(name='doc_id', dtype=DataType.INT64,
                                     description='doc_id', max_length=500,
                                     is_primary=True, auto_id=False)
                title_vec = FieldSchema(name='title_vec', dtype=DataType.FLOAT_VECTOR,
                                        description='title embedding vectors',
                                        dim=VECTOR_DIMENSION, is_primary=False)
                body_vec = FieldSchema(name='body_vec', dtype=DataType.FLOAT_VECTOR,
                                       description='body embedding vectors',
                                       dim=VECTOR_DIMENSION, is_primary=False)

                schema = CollectionSchema(fields=[doc_id, title_vec, body_vec], description='vec_info')
                self.collection = Collection(name=collection_name, schema=schema)
                self.create_collection(collection_name)
                LOGGER.info(f"create milvus collection:{collection_name}")
            else:
                self.set_collection(collection_name)
            return "ok"
        except Exception as e:
            LOGGER.error(f"failed to create collection {e}")
            sys.exit(1)

    def insert(self, collection_name, doc_id, title_vec, body_vec):
        try:
            data = [collection_name, doc_id, title_vec, body_vec]
            self.set_collection(collection_name)
            mr = self.collection.insert(data)
            ids = mr.primary_keys
            self.collection.load()
            LOGGER.info(
                f"insert vector to milvus in collection: {collection_name} with body {len(body_vec)}, {len(title_vec)}")
            return ids
        except Exception as e:
            LOGGER.error(f"insert to load data to milvus {e}")
            sys.exit(1)

    def create_index(self, collection_name):
        try:
            self.set_collection(collection_name)
            default_index = {"index_type": "IVF_SQ8", "params": {"nlist": 16384}}
            status = self.collection.create_index(field_name="embedding", index_params=default_index)
            if not status:
                LOGGER.info(f"successfully create index in collection :{collection_name} with param:{default_index}")
                return status
            else:
                raise Exception(status.message)
        except Exception as e:
            LOGGER.error(f"failed to create index:{e}")
            sys.exit(1)

    def delete_collection(self, collection_name):
        try:
            self.set_collection(collection_name)
            self.collection.drop()
            LOGGER.info("successfully drop collection!")
            return "ok"
        except Exception as e:
            LOGGER.error(f"failed to drop collection :{e}")
            sys.exit(1)

    def search_vector(self, collection_name, vectors, top_k):
        try:
            self.set_collection(collection_name)
            search_params = {"params": {"nprobe": 16}}
            res = self.collection.search(vectors, anns_field="embeding", param=search_params, limit=top_k)
            LOGGER.info(f"successfully search in collection:{res}", res)
            return res
        except Exception as e:
            LOGGER.error(f"failed to search vectors in milvus:{e}")
            sys.exit(1)

    def count(self, collection_name):
        try:
            self.set_collection(collection_name)
            num = self.collection.num_entities
            LOGGER.info(f"successfully get the num:{num} of the collection:{collection_name}")
            return num
        except Exception as e:
            LOGGER.error(f"failed to count vectors in milvus:{e}")
            sys.exit(1)
