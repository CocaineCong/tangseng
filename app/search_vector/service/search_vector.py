"""search vector grpc service"""
import json
import grpc
import logging
import asyncio

from ..consts.consts import VECTOR_RECALL_TOPK
from idl.pb.search_vector import search_vector_pb2
from ..config.config import DEFAULT_MILVUS_TABLE_NAME, VECTOR_ADDR
from ..milvus.operators import do_search
from ..etcd_operate.etcd import etcd_client
from ..milvus.milvus import milvus_client
from idl.pb.search_vector import search_vector_pb2_grpc


class SearchVectorService(search_vector_pb2_grpc.SearchVectorServiceServicer):
    """
    search vector service objective
    """

    def SearchVector(self, request,
                     context) -> search_vector_pb2.SearchVectorResponse:
        try:
            queryies = request.query
            doc_ids = []
            for query in queryies:
                ids, distants = do_search(DEFAULT_MILVUS_TABLE_NAME, query,
                                          VECTOR_RECALL_TOPK, milvus_client)
                print("search vector ids", ids)
                doc_ids += ids
            print("search vector data", doc_ids)
            return search_vector_pb2.SearchVectorResponse(code=200,
                                                          doc_ids=doc_ids,
                                                          msg='ok',
                                                          error='')
        except Exception as e:
            print("search vector error", e)
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
    await server.start()
    await server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(serve())
