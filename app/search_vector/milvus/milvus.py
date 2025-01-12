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

import sys
from pymilvus import (connections, Collection, utility, FieldSchema, DataType,
                      CollectionSchema)
from ..config.config import VECTOR_DIMENSION, MILVUS_HOST, MILVUS_PORT, METRIC_TYPE
from ..utils.logs import LOGGER


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
                doc_id = FieldSchema(name='doc_id',
                                     dtype=DataType.INT64,
                                     description='doc_id',
                                     max_length=500,
                                     is_primary=True,
                                     auto_id=False)
                # One and only vector can be created.
                # tow or zero will be wrong.
                body_vec = FieldSchema(name='embedding',
                                       dtype=DataType.FLOAT_VECTOR,
                                       description='body embedding vectors',
                                       dim=VECTOR_DIMENSION,
                                       is_primary=False)

                schema = CollectionSchema(fields=[doc_id, body_vec],
                                          description='vec_info')
                self.collection = Collection(name=collection_name,
                                             schema=schema)
                self.create_index(collection_name)
                LOGGER.info(f"create milvus collection:{collection_name}")
            else:
                self.set_collection(collection_name)
            return "ok"
        except Exception as e:
            LOGGER.error(f"failed to create collection {e}")
            sys.exit(1)

    def insert(self, collection_name, doc_id, body_vec):
        try:
            data = [doc_id, body_vec]
            self.set_collection(collection_name)
            mr = self.collection.insert(data)
            ids = mr.primary_keys
            self.collection.load()
            LOGGER.info(
                f"insert vector to milvus in collection: {collection_name} with body {len(body_vec)}"
            )
            return ids
        except Exception as e:
            LOGGER.error(f"insert to load data to milvus {e}")
            sys.exit(1)

    def create_index(self, collection_name):
        try:
            self.set_collection(collection_name)
            default_index = {
                "index_type": "IVF_SQ8",
                "metric_type": METRIC_TYPE,
                "params": {
                    "nlist": 16384
                }
            }
            print(default_index)
            status = self.collection.create_index(field_name="embedding",
                                                  index_params=default_index)
            if not status.code:
                LOGGER.info(
                    f"successfully create index in collection :{collection_name} with param:{default_index}"
                )
                return status
            else:
                raise Exception(status.message)
        except Exception as e:
            print(e)
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

    def search_vectors(self, collection_name, vectors, top_k):
        try:
            self.set_collection(collection_name)
            search_params = {
                "metric_type": METRIC_TYPE,
                "params": {
                    "nprobe": 16
                }
            }
            res = self.collection.search(vectors,
                                         anns_field="embedding",
                                         param=search_params,
                                         limit=top_k)
            LOGGER.info(f"successfully search in collection:{res}")
            return res
        except Exception as e:
            LOGGER.error(f"failed to search vectors in milvus:{e}")
            sys.exit(1)

    def count(self, collection_name):
        try:
            self.set_collection(collection_name)
            num = self.collection.num_entities
            LOGGER.info(
                f"successfully get the num:{num} of the collection:{collection_name}"
            )
            return num
        except Exception as e:
            LOGGER.error(f"failed to count vectors in milvus:{e}")
            sys.exit(1)


milvus_client = Milvus()
