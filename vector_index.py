"""the script file is to handle vector index from kafka"""
import threading
from app.search_vector.consts.consts import KAFKA_CONSUMER_VECTOR_INDEX_TOPIC
from app.search_vector.config.config import DEFAULT_MILVUS_TABLE_NAME
from app.search_vector.kafka_operate.consumer import store_data_from_kafka


def consume_inverted_index():
    """
    consume data from kafka to build inverted index
    """
    topic = KAFKA_CONSUMER_VECTOR_INDEX_TOPIC
    table_name = DEFAULT_MILVUS_TABLE_NAME
    thread = threading.Thread(target=store_data_from_kafka(
        topic, table_name))  # 创建线程对象
    thread.start()  # 启动线程
    print("start consume inverted index")
    thread.join()  # 等待线程结束


if __name__ == "__main__":
    consume_inverted_index()
