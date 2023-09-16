package test

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/codec"
	"github.com/CocaineCong/tangseng/types"
)

func TestStorageInverted(t *testing.T) {
	// 读取文件
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	// 建立倒排索引
	p := &types.InvertedIndexValue{
		Token:         "测试文本",
		PostingsList:  segment.CreateNewPostingsList(1),
		DocCount:      10,
		PositionCount: 20,
		TermValues:    nil,
	}
	// 编码
	buf, err := codec.EncodePostings(p)
	if err != nil {
		log.LogrusObj.Errorf("updatePostings encodePostings err: %v", err)
		return
	}
	// 倒排索引表
	inverted := storage.NewInvertedDB(termName, postingsName)
	err = inverted.StoragePostings(p.Token, buf)
	if err != nil {
		fmt.Println(err)
	}
}

func TestStoreInverted(t *testing.T) {
	// 读取文件
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	// 建立倒排索引
	p := &types.InvertedIndexValue{
		Token:         "测试文本",
		PostingsList:  segment.CreateNewPostingsList(1),
		DocCount:      10,
		PositionCount: 20,
		TermValues: &types.TermValue{
			DocCount: 10,
			Offset:   1,
			Size:     2,
		},
	}
	// 编码
	buf, err := codec.EncodePostings(p)
	if err != nil {
		fmt.Println(err)
	}
	// 倒排索引表
	inverted := storage.NewInvertedDB(termName, postingsName)
	err = inverted.StoragePostings(p.Token, buf)
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
	invertedValue, err := inverted.GetInvertedInfo(token)
	if err != nil {
		fmt.Println(err)
	}
	// 编码
	fmt.Println(invertedValue)
}
