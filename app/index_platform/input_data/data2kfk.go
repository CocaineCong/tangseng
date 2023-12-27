package input_data

import (
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	"github.com/CocaineCong/tangseng/types"
	"github.com/pkg/errors"
)

// DocData2Kfk Doc数据处理
func DocData2Kfk(doc *types.Document) (err error) {
	doctByte, _ := doc.MarshalJSON()
	err = kfk.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		return errors.WithMessagef(err, "DocData2Kfk-KafkaCSVLoaderTopic :%v", err)
	}

	return
}

// DocTrie2Kfk Trie树构建
func DocTrie2Kfk(tokens []string) (err error) {
	for _, k := range tokens {
		err = kfk.KafkaProducer(consts.KafkaTrieTreeTopic, []byte(k))
	}

	if err != nil {
		return errors.WithMessagef(err, "DocTrie2Kfk-KafkaTrieTreeTopic :%v", err)
	}

	return
}
