from kafka import KafkaProducer, KafkaConsumer
from kafka.errors import KafkaError
from milvus import milvus
from milvus.operators import do_upload


class KafkaHelper:

    def __init__(self, bootstrap_servers):
        self.bootstrap_servers = bootstrap_servers
        self.producer = None
        self.consumer = None
        self.milvus_client = milvus.milvus_client

    def connect_producer(self):
        try:
            self.producer = KafkaProducer(bootstrap_servers=self.bootstrap_servers)
            print("Connected to Kafka producer successfully.")
        except KafkaError as e:
            print(f"Failed to connect to Kafka producer: {e}")

    def connect_consumer(self, topic):
        try:
            self.consumer = KafkaConsumer(topic, bootstrap_servers=self.bootstrap_servers)
            print(f"Connected to Kafka consumer successfully. Listening to topic: {topic}")
        except KafkaError as e:
            print(f"Failed to connect to Kafka consumer: {e}")

    def send_message(self, topic, message):
        if not self.producer:
            self.connect_producer()

        try:
            self.producer.send(topic, message.encode('utf-8')).add_callback(self.on_send_success).add_errback(self.on_send_error)
            self.producer.flush()
        except KafkaError as e:
            print(f"Failed to send message to Kafka: {e}")

    def consume_messages(self):
        if not self.consumer:
            print("No Kafka consumer connected.")
            return
        print("Consuming messages...")
        for message in self.consumer:
            # print(message.value.decode('utf-8'))
            do_upload("test", 1, "mirror",
                    message, self.milvus_client)

            
    def on_send_success(self, record_metadata):
        print(f"Message sent successfully. Topic: {record_metadata.topic}, Partition: {record_metadata.partition}, Offset: {record_metadata.offset}")

    def on_send_error(self, exception):
        print(f"Failed to send message to Kafka: {exception}")



if __name__ == "__main__":
    kafka_servers = ['localhost:10001','localhost:10002','localhost:10003']
    kafka_topic = 'my_topic'

    kafka_helper = KafkaHelper(kafka_servers)
    kafka_helper.connect_consumer(kafka_topic)

    while True:
        message = input("Enter a message to send (or 'q' to quit): ")
        if message.lower() == 'q':
            break
        kafka_helper.send_message(kafka_topic, message)
