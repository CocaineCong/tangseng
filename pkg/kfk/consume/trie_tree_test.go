package consume

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	storage.InitTrieDBs()
	// trie.InitTrieTree()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestTrieTreeConsumerPutTrieTree(t *testing.T) {
	key := []byte(consts.TrieTreeBucket)
	value := []byte("test")
	err := storage.GlobalTrieDBs.PutTrieTree(key, value)
	if err != nil {
		fmt.Println(err)
	}
}

func TestTrieTreeKafkaGetTrieTree(t *testing.T) {
	key := []byte(consts.TrieTreeBucket)
	value, err := storage.GlobalTrieDBs.GetTrieTree(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value))
}
