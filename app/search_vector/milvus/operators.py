import sys

from ..config.config import DEFAULT_MILVUS_TABLE_NAME, TRANSFORMER_MODEL
from ..utils.logs import LOGGER


def do_upload(table_name, doc_id, title, body, milvus_client):
    """
    upload info in milvus

    :params table_name the table name of milvus
    :params doc_id the id of doc
    :params title the title of doc
    :params body the body of doc
    :params milvus_client milvus client
    """
    try:
        if not table_name:
            table_name = DEFAULT_MILVUS_TABLE_NAME
        milvus_client.create_collection(table_name)
        body_feat = TRANSFORMER_MODEL.encode(title + body)  # word 转 vec
        ids = milvus_client.insert(table_name, [doc_id], [body_feat])
        return ids
    except Exception as e:
        LOGGER.error(f"failed with upload :{e}")
        sys.exit(1)


def do_search(table_name, query, top_k, milvus_client):
    """
    query 传入的搜索参数, 返回 doc ids 及其 距离参数
    
    :params table_name the table name of milvus
    :params query the query of search
    :params top_k the top k of search
    :params milvus_client milvus client
    
    :return doc_ids the doc ids of search
    :return distances the distances of search
    """
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
