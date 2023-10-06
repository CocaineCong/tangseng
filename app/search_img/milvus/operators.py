import sys

from config.config import DEFAULT_MILVUS_TABLE_NAME
from utils.encode import word2vec
from utils.logs import LOGGER


# 上传数据到milvus中
def do_upload(table_name, doc_id, title, body, milvus_client):
    try:
        if not table_name:
            table_name = DEFAULT_MILVUS_TABLE_NAME
        milvus_client.create_collection(table_name)
        title_feat = word2vec(title)  # word 转 vec
        body_feat = word2vec(body)
        ids = milvus_client.insert(table_name, doc_id, [title_feat], [body_feat])
        return ids
    except Exception as e:
        LOGGER.error(f"failed with upload :{e}")
        sys.exit(1)


# query 传入的搜索参数
def do_search(table_name, query, top_k, milvus_client):
    try:
        if not table_name:
            table_name = DEFAULT_MILVUS_TABLE_NAME
        query_feat = word2vec(query)
        vectors = milvus_client.search_vectors(table_name, [query_feat], top_k)
        doc_ids = [str(x.id) for x in vectors[0]]
        distances = [x.distance for x in vectors[0]]
        return doc_ids, distances
    except Exception as e:
        LOGGER.error(f"error with search : {e}")
        sys.exit(1)
