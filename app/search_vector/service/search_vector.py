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
"""search vector grpc service"""
import inspect
import json
import grpc
import logging
import asyncio

from opentelemetry import trace
from ..consts.consts import VECTOR_RECALL_TOPK
from idl.pb.search_vector import search_vector_pb2
from ..config.config import DEFAULT_MILVUS_TABLE_NAME, VECTOR_ADDR
from ..milvus.operators import do_search
from ..etcd_operate.etcd import etcd_client
from ..milvus.milvus import milvus_client
from idl.pb.search_vector import search_vector_pb2_grpc

tracer = trace.get_tracer(__name__)


class SearchVectorService(search_vector_pb2_grpc.SearchVectorServiceServicer):
    """
    search vector service objective
    """

    def SearchVector(self, request,
                     context) -> search_vector_pb2.SearchVectorResponse:
        with tracer.start_as_current_span(
                inspect.getframeinfo(inspect.currentframe()).function):
            try:
                queries = request.query
                doc_ids = []
                for query in queries:
                    ids, distant = do_search(DEFAULT_MILVUS_TABLE_NAME, query,
                                             VECTOR_RECALL_TOPK,
                                             milvus_client)
                    doc_ids += ids
                return search_vector_pb2.SearchVectorResponse(code=200,
                                                              doc_ids=doc_ids,
                                                              msg='ok',
                                                              error='')
            except Exception as e:
                return search_vector_pb2.SearchVectorResponse(code=500,
                                                              error=str(e))


async def serve() -> None:
    server = grpc.aio.server()
    search_vector_pb2_grpc.add_SearchVectorServiceServicer_to_server(
        SearchVectorService(), server)
    server.add_insecure_port(VECTOR_ADDR)
    logging.info("Starting server on %s", VECTOR_ADDR)
    key = f"/search_vector/{VECTOR_ADDR}"
    value = json.dumps({
        "name": "search_vector",
        "addr": f"{VECTOR_ADDR}",
        "version": "",
        "weight": 0
    })
    etcd_client.set(key, value)
    logging.info("set %s node on %s", key, value)
    await server.start()
    await server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(serve())
