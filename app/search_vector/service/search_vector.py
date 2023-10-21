import grpc
import logging
import asyncio

from idl.pb import search_vector_pb2_grpc
from ctl.resp import RespDefaultSuccess,RespDefaultError
from config.config import DEFAULT_MILVUS_TABLE_NAME
from milvus.operators import do_search
from milvus.milvus import milvus_client

class SearchVectorService(search_vector_pb2_grpc.SearchVectorServiceServicer):
    def SearchVector(self, request, context):
        try:
            queryies = request.query
            doc_ids = []
            for query in queryies:
                ids, distants = do_search(DEFAULT_MILVUS_TABLE_NAME, query, 1, milvus_client)
                doc_ids.append(ids)
            data = ','.join(doc_ids)
            return RespDefaultSuccess(data)
        except Exception as e:
            return RespDefaultError(str(e))

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
