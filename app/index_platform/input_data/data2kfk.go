package input_data

import (
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// DocData2Kfk Doc数据处理
func DocData2Kfk(doc *types.Document) (err error) {
	doctByte, _ := doc.MarshalJSON()
	err = kfk.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		logs.LogrusObj.Errorf("DocData2Kfk-KafkaCSVLoaderTopic :%+v", err)
		return
	}

	return
}

// DocTrie2Kfk Trie树构建
func DocTrie2Kfk(tokens []string) (err error) {
	for _, k := range tokens {
		err = kfk.KafkaProducer(consts.KafkaTrieTreeTopic, []byte(k))
	}

	if err != nil {
		logs.LogrusObj.Errorf("DocTrie2Kfk-KafkaTrieTreeTopic :%+v", err)
		return
	}

	return
}
