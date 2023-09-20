package test

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/config"
)

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
	fmt.Println(invertedValue)
}
