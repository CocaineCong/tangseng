package storage

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
)

func TestInvertedDBRead(t *testing.T) {
	query := "导演"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	termValue, err := inverted.GetTermInfo(query)
	if err != nil {
		fmt.Println("Err", err)
	}
	fmt.Println("termValue", termValue)
	v, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v", v) // TODO:存储有问题
	err = inverted.StoragePostings(query, []byte("100"), 1)
	v2, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v2", string(v2))
}
