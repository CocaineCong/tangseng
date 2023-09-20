package storage

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
)

func TestInvertedDBRead(t *testing.T) {
	query := "电影"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	termValue, err := inverted.GetTermInfo(query)
	if err != nil {
		fmt.Println("Err", err)
	}
	fmt.Println("termValue", termValue)
}

func TestGetInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	p, err := inverted.GetInvertedInfo(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Token, p.PostingsList)
}

func TestInitInvertedDB(t *testing.T) {
	InitInvertedDB()
}
