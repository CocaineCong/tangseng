import sys

from config.config import DEFAULT_MILVUS_TABLE_NAME, TRANSFORMER_MODEL
from utils.logs import LOGGER


# 上传数据到milvus中
def do_upload(table_name, doc_id, title, body, milvus_client):
    try:
        if not table_name:
            table_name = DEFAULT_MILVUS_TABLE_NAME
        # if milvus_client.has_collection(table_name):
        #     milvus_client.delete_collection(table_name)
        milvus_client.create_collection(table_name)
        body_feat = TRANSFORMER_MODEL.encode(title+body)  # word 转 vec
        ids = milvus_client.insert(table_name, [doc_id], [body_feat])
        return ids
    except Exception as e:
        LOGGER.error(f"failed with upload :{e}")
        sys.exit(1)


# query 传入的搜索参数, 返回 doc ids 及其 距离参数
def do_search(table_name, query, top_k, milvus_client):
    try:
        if not table_name:
            table_name = DEFAULT_MILVUS_TABLE_NAME
        query_feat = TRANSFORMER_MODEL.encode(query)
        vectors = milvus_client.search_vectors(table_name, [query_feat], top_k)
        doc_ids = [str(x.id) for x in vectors[0]]
        distances = [x.distance for x in vectors[0]]
        return doc_ids, distances
    except Exception as e:
        LOGGER.error(f"error with search : {e}")
        sys.exit(1)
