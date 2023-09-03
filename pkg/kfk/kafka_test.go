package kfk

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../config/config.yaml"}
	config.InitConfigForTest(&re)
	InitKafka()
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

type TestKafkaData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TestKafkaProducer(t *testing.T) {
	data := &TestKafkaData{
		Key:   "怎么说",
		Value: "啊哈哈哈哈",
	}
	d, _ := json.Marshal(data)
	err := KafkaProducer("search-engine-csv-loader-topic", string(d))
	if err != nil {
		fmt.Println("KafkaProducer", err)
		return
	}
	fmt.Println("Produce Message Finish")
}

func TestKafkaConsumer(t *testing.T) {
	msgs, err := KafkaConsumer("search-engine-csv-loader-topic")
	if err != nil {
		fmt.Println("KafkaProducer", err)
		return
	}
	for {
		select {
		case msg := <-msgs:
			key := string(msg.Key)
			value := string(msg.Value)
			fmt.Println(key, value)
		}
	}

}

func TestKafkaProduce(t *testing.T) {
	KafkaProduce()
}

func TestKafkaConsume(t *testing.T) {
	KafkaConsume()
}
