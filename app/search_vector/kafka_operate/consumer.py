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

"""store vector index from kafka"""
import inspect
from opentelemetry import trace
from ..kafka_operate.kafka_operate import kafka_helper


def store_data_from_kafka(kafka_topic, milvus_table_name):
    """
    store data to milvus from kafka for building inverted index
    """
    tracer = trace.get_tracer(__name__)
    with tracer.start_as_current_span(inspect.getframeinfo(inspect.currentframe()).function):
        kafka_helper.connect_consumer(kafka_topic)
        kafka_helper.consume_messages_store_milvus(milvus_table_name)
