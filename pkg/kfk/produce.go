// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package kfk

import (
	_ "net/http/pprof"

	"github.com/pkg/errors"

	"github.com/IBM/sarama"
)

// KafkaProducer 发送单条
func KafkaProducer(topic string, msg []byte) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(GobalKafka)
	if err != nil {
		return errors.Wrap(err, "failed to create Kafka producer")
	}
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	_, _, err = producer.SendMessage(message)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}
	return
}

// KafkaProducers 发送多条，topic在messages中
func KafkaProducers(messages []*sarama.ProducerMessage) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(GobalKafka)
	if err != nil {
		return errors.Wrap(err, "failed to create Kafka producer")
	}
	err = producer.SendMessages(messages)
	if err != nil {
		return errors.Wrap(err, "failed to send messages")
	}
	return
}
