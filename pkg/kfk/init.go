package kfk

import (
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/config"
)

var GobalKafka sarama.Client

func InitKafka() {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, con)
	if err != nil {
		logs.LogrusObj.Errorln(err)
		return
	}
	GobalKafka = kafkaClient
}
