package input_data

import (
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// DocData2Kfk Doc数据处理
func DocData2Kfk(doc *types.Document) {
	doctByte, _ := doc.MarshalJSON()
	err := kfk.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		logs.LogrusObj.Errorf("DocData2Kfk-KafkaCSVLoaderTopic :%+v", err)
	}
}

// DocTrie2Kfk Trie树构建
func DocTrie2Kfk(token string) {
	err := kfk.KafkaProducer(consts.KafkaTrieTreeTopic, []byte(token))
	if err != nil {
		logs.LogrusObj.Errorf("DocTrie2Kfk-KafkaTrieTreeTopic :%+v", err)
	}
}
