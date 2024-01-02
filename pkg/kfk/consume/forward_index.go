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
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"

	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/app/index_platform/repository/db/dao"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/model"
	"github.com/CocaineCong/tangseng/types"
)

// ForwardIndexKafkaConsume 正排索引的消费建立
func ForwardIndexKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {
	keepRunning := true
	logs.LogrusObj.Infof("Starting a new Sarama consumer")
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	// 设置一个消费组
	consumer := ForwardIndexConsumer{
		Ready: make(chan bool),
	}
	configK := kfk.GetDefaultConsumeConfig(assignor)
	cancelCtx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, configK)
	if err != nil {
		logs.LogrusObj.Errorf("Error creating consumer group woker: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	logs.LogrusObj.Infof("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-cancelCtx.Done():
			logs.LogrusObj.Infof("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			logs.LogrusObj.Infof("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		err = errors.Wrapf(err, "failed to close woker")
		return
	}

	return
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		logs.LogrusObj.Infof("Resuming consumption")
	} else {
		client.PauseAll()
		logs.LogrusObj.Infof("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer Sarama消费者群体的消费者
type ForwardIndexConsumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ForwardIndexConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ForwardIndexConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 必须启动 ConsumerGroupClaim 的 Messages() 消费者循环。
// 一旦 Messages() 通道关闭，处理程序必须完成其处理循环并退出。
func (consumer *ForwardIndexConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()
	task := &types.Task{
		Columns:    []string{"doc_id", "title", "body", "url"},
		BiTable:    "data",
		SourceType: consts.DataSourceCSV,
	}
	iDao := dao.NewInputDataDao(ctx)
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logs.LogrusObj.Infof("message channel was closed")
				return nil
			}

			if task.SourceType == consts.DataSourceCSV {
				doc := new(types.Document)
				_ = doc.UnmarshalJSON(message.Value)
				// TODO: 后续再开发starrocks
				// up.Push(&types.Data2Starrocks{
				// 	DocId: docs.DocId,
				// 	Url:   "",
				// 	Title: docs.Title,
				// 	Desc:  docs.Body,
				// 	Score: 0.00, // 评分
				// })
				inputData := &model.InputData{
					DocId:  doc.DocId,
					Title:  doc.Title,
					Body:   doc.Body,
					Url:    "",
					Score:  0.0,
					Source: task.SourceType,
				}
				err := iDao.CreateInputData(inputData)
				if err != nil {
					logs.LogrusObj.Errorf("iDao.CreateInputData:%+v", err)
				}
			}

			logs.LogrusObj.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
