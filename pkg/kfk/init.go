package kfk

import (
	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/config"
)

var GobalKafka sarama.Client

func InitKafka() {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, con)
	if err != nil {
		return
	}
	GobalKafka = kafkaClient
}
