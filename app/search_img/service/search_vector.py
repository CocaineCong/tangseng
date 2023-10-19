import grpc
import logging
import asyncio

from app.search_img.config.config import DEFAULT_MILVUS_TABLE_NAME
from app.search_img.milvus.operators import do_search
from app.search_img.milvus.milvus import milvus_client
from ....idl.pb.search_vector import search_vector_pb2,search_vector_pb2_grpc

class SearchVectorService(search_vector_pb2_grpc.SearchVectorServiceServicer):

    def SearchVector(self, request, context):
        resp = {'code':200, 'doc_ids':[], 'msg':""}
        try:
            queryies = request.query
            doc_ids = []
            for query in queryies:
                ids,distants = do_search(DEFAULT_MILVUS_TABLE_NAME, query, 1, milvus_client)
                doc_ids.append(ids)
            resp['doc_ids'] = doc_ids
        except Exception as e:
            resp['code'] = 500
            resp['msg'] = str(e)

        return resp


async def serve() -> None:
    server = grpc.aio.server()
    search_vector_pb2_grpc.add_SearchVectorServiceServicer_to_server(SearchVectorService(), server)
    listen_addr = "[::]:50051"
    server.add_insecure_port(listen_addr)
    logging.info("Starting server on %s", listen_addr)
    await server.start()
    await server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(serve())
