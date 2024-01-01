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

package consume

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/app/index_platform/trie"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

// TrieTreeKafkaConsume token词的消费建立
func TrieTreeKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {
	logs.LogrusObj.Infof("Starting a new Sarama consumer")
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// 设置一个消费组
	consumer := TrieTreeConsumer{
		Ready: make(chan bool),
	}

	configK := kfk.GetDefaultConsumeConfig(assignor)
	cancelCtx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, configK)
	if err != nil {
		logs.LogrusObj.Errorf("Error creating consumer group woker: %v", err)
	}

	go func() {
		for {
			if err = client.Consume(cancelCtx, []string{topic}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				logs.LogrusObj.Errorf("Error from consumer: %v", err)
			}
			if cancelCtx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready
	cancel()

	return
}

// TrieTreeConsumer Sarama消费者群体的消费者
type TrieTreeConsumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *TrieTreeConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *TrieTreeConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 必须启动 ConsumerGroupClaim 的 Messages() 消费者循环。
// 一旦 Messages() 通道关闭，处理程序必须完成其处理循环并退出。
func (consumer *TrieTreeConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// ctx := context.Background()
	gapTime := 2 * time.Minute
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logs.LogrusObj.Infof("message channel was closed")
				return nil
			}
			// 构建trie tree树
			trie.GobalTrieTree.Insert(string(message.Value))
			// logs.LogrusObj.Infof("TrieTreeConsumer Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// https://github.com/IBM/sarama/issues/1192

		case <-time.After(gapTime):
			logs.LogrusObj.Infof("ConsumeClaim starting store dict")
			// _ = storage.GlobalTrieDBs.StorageDict(trie.GobalTrieTree) // TODO:后续看看能不能实现一个全局的triedb，每次都先读取存量进行初始化，再插入增量...
			logs.LogrusObj.Infof("ConsumeClaim ending store dict")

		case <-session.Context().Done():
			logs.LogrusObj.Infof("TrieTreeConsumer Done!")
			return nil
		}
	}
}

// func mergeTrieTree(node string) {
// 	trie.GobalTrieTree.Insert(node)
// 	gapTime := 2 * time.Minute
// 	for {
// 		select {
// 		case <-time.After(gapTime):
// 			_ = storage.GlobalTrieDBs.StorageDict(trie.GobalTrieTree)
// 		}
// 	}
// }
