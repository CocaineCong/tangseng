package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
)

func TestGetInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	p, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}

func TestInitInvertedDB(t *testing.T) {
	ctx := context.Background()
	InitInvertedDB(ctx)
	for _, v := range GlobalInvertedDB {
		fmt.Println(v)
	}
}
