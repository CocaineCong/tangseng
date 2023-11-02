from kafka_operate import KafkaHelper


if __name__ == "__main__":
    kafka_servers = ['localhost:10001','localhost:10002','localhost:10003']
    kafka_topic = 'my_topic'

    kafka_helper = KafkaHelper(kafka_servers)
    kafka_helper.connect_consumer(kafka_topic)

    kafka_helper.consume_messages()
