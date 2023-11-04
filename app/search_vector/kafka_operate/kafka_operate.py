"""kafka operate"""
import json
from kafka import KafkaProducer, KafkaConsumer
from kafka.errors import KafkaError
from ..config.config import KAFKA_CLUSTER
from ..milvus import milvus
from ..milvus.operators import do_upload


class KafkaHelper:
    """
    kafka handle objective
    """

    def __init__(self, bootstrap_servers):
        self.bootstrap_servers = bootstrap_servers
        self.producer = None
        self.consumer = None
        self.milvus_client = milvus.milvus_client

    def connect_producer(self):
        """
        connect kafka producer
        """
        try:
            self.producer = KafkaProducer(
                bootstrap_servers=self.bootstrap_servers)
            print("Connected to Kafka producer successfully.")
        except KafkaError as e:
            print(f"Failed to connect to Kafka producer: {e}")

    def connect_consumer(self, topic):
        """
        connect kafka consumer
        """
        try:
            self.consumer = KafkaConsumer(
                topic, bootstrap_servers=self.bootstrap_servers)
            print(
                f"Connected to Kafka consumer successfully. Listening to topic: {topic}"
            )
        except KafkaError as e:
            print(f"Failed to connect to Kafka consumer: {e}")

    def send_message(self, topic, msg):
        """
        send message by topic
        """
        if not self.producer:
            self.connect_producer()

        try:
            self.producer.send(topic, msg.encode('utf-8')).add_callback(
                self.on_send_success).add_errback(self.on_send_error)
            self.producer.flush()
        except KafkaError as e:
            print(f"Failed to send message to Kafka: {e}")

    def consume_messages(self):
        """
        consume messages from kafka
        """
        if not self.consumer:
            print("No Kafka consumer connected.")
            return
        print("Consuming messages...")
        for msg in self.consumer:
            print(msg)

    def consume_messages_store_milvus(self, milvus_table):
        """
        consume messages from kafka and store in milvus
        """
        if not self.consumer:
            print("No Kafka consumer connected.")
            return
        print("Consuming messages...")
        for msg in self.consumer:
            data = json.loads(msg.value.decode('utf-8'))
            do_upload(milvus_table, int(data["doc_id"]), data["title"],
                      data["body"], self.milvus_client)

    def on_send_success(self, record_metadata):
        """
        a hook when send message to kafka is success
        """
        print(f"Message sent successfully. Topic: {record_metadata.topic}")
        print(
            f"Partition: {record_metadata.partition},Offset: {record_metadata.offset}"
        )

    def on_send_error(self, exception):
        """
        a hook when send message to kafka is failed
        """
        print(f"Failed to send message to Kafka: {exception}")


kafka_helper = KafkaHelper(KAFKA_CLUSTER)
