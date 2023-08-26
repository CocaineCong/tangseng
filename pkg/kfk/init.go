package kfk

import (
	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/config"
)

var GobalKafka sarama.Client

func InitKafka() {
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, nil)
	if err != nil {
		return
	}
	GobalKafka = kafkaClient
}
