from kafka_operate import KafkaHelper
import threading


def consumer_inverted_index():
    kafka_servers = ['localhost:10001','localhost:10002','localhost:10003']
    kafka_topic = 'my_topic'
    kafka_helper = KafkaHelper(kafka_servers)
    kafka_helper.connect_consumer(kafka_topic)
    kafka_helper.consume_messages()

# 创建线程对象
thread = threading.Thread(target=consumer_inverted_index)
# 启动线程
thread.start()
# 等待线程结束
thread.join()
