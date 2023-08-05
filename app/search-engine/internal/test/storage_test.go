package test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/segment"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestStorageInverted(t *testing.T) {
	// 读取文件
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	// 建立倒排索引
	p := &segment.InvertedIndexValue{
		Token:         "测试文本",
		PostingsList:  segment.CreateNewPostingsList(1),
		DocCount:      10,
		PositionCount: 20,
		TermValues:    nil,
	}
	// 编码
	buf, err := segment.EncodePostings(p.PostingsList, p.DocCount)
	if err != nil {
		log.LogrusObj.Errorf("updatePostings encodePostings err: %v", err)
		return
	}
	// 倒排索引表
	inverted := storage.NewInvertedDB(termName, postingsName)
	err = inverted.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
	if err != nil {
		fmt.Println(err)
	}
}

func TestStoreInverted(t *testing.T) {
	// 读取文件
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	// 建立倒排索引
	p := &segment.InvertedIndexValue{
		Token:         "测试文本",
		PostingsList:  segment.CreateNewPostingsList(1),
		DocCount:      10,
		PositionCount: 20,
		TermValues:    nil,
	}
	// 编码
	buf, err := segment.EncodePostings(p.PostingsList, p.DocCount)
	if err != nil {
		fmt.Println(err)
	}
	// 倒排索引表
	inverted := storage.NewInvertedDB(termName, postingsName)
	err = inverted.StoragePostings(p.Token, buf.Bytes(), p.DocCount)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetInverted(t *testing.T) {
	// 读取文件
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	token := "测试文本"
	inverted := storage.NewInvertedDB(termName, postingsName)
	invertedValue, err := inverted.GetInverted([]byte(token))
	if err != nil {
		fmt.Println(err)
	}
	// 编码
	p, c, err := segment.DecodePostings(bytes.NewBuffer(invertedValue))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p, c)
}
