# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

"""the script file is to handle vector index from kafka"""
import threading
from opentelemetry import trace
from app.search_vector.consts.consts import KAFKA_CONSUMER_VECTOR_INDEX_TOPIC
from app.search_vector.config.config import DEFAULT_MILVUS_TABLE_NAME
from app.search_vector.kafka_operate.consumer import store_data_from_kafka


def consume_inverted_index():
    """
    consume data from kafka to build inverted index
    """
    topic = KAFKA_CONSUMER_VECTOR_INDEX_TOPIC
    table_name = DEFAULT_MILVUS_TABLE_NAME
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span("consume_inverted_index"):
        thread = threading.Thread(target=store_data_from_kafka(
            topic, table_name))  # 创建线程对象
        thread.start()  # 启动线程
        print("start consume inverted index")
        thread.join()  # 等待线程结束


if __name__ == "__main__":
    consume_inverted_index()
