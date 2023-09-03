package kfk

import (
	"github.com/IBM/sarama"
)

func KafkaConsumer(topic string) (msg <-chan *sarama.ConsumerMessage, err error) {
	consumer, err := sarama.NewConsumerFromClient(GobalKafka)
	if err != nil {
		return
	}
	var partition int32 = 0
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return
	}

	return partitionConsumer.Messages(), nil
}
