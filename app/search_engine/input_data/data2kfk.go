package input_data

import (
	"github.com/IBM/sarama"

	"github.com/samber/lo"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const inputDataPath = "./movies_data.csv"

// DocData2Kfk Doc数据处理
func DocData2Kfk() {
	docs := ReadFiles([]string{inputDataPath})
	data2kfkList := make([]*sarama.ProducerMessage, 0)
	for _, doc := range docs[1:] {
		doct, _ := doc2Struct(doc)
		doctByte, _ := doct.MarshalJSON()
		data2kfkList = append(data2kfkList, &sarama.ProducerMessage{
			Topic: consts.KafkaCSVLoaderTopic,
			Key:   nil,
			Value: sarama.StringEncoder(doctByte),
		})
	}

	// 200 一砸生产
	producers := lo.Chunk(data2kfkList, consts.KafkaBatchProduceCount)
	for _, producer := range producers {
		err := kfk.KafkaProducers(producer)
		if err != nil {
			logs.LogrusObj.Errorf("DocData2Kfk-KafkaProducers :%+v", err)
			return
		}
	}
}
