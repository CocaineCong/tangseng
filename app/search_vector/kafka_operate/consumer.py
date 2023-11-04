"""store vector index from kafka"""
from kafka_operate.kafka_operate import kafka_helper


def store_data_from_kafka(kafka_topic, milvus_table_name):
    """
    store data to mivlus from kakfa for building inverted index
    """
    kafka_helper.connect_consumer(kafka_topic)
    kafka_helper.consume_messages_store_milvus(milvus_table_name)
