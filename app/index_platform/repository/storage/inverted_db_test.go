package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
)

func TestInvertedDBRead(t *testing.T) {
	query := "电影"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	inverted := NewInvertedDB(termName)
	v, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v", string(v))
	err = inverted.StoragePostings(query, []byte("100"))
	v2, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v2", string(v2))

	inverted.PutInverted([]byte(query), []byte("11111"))
	v3, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(v3))
}

func TestStoreInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	output := roaring.New()
	output.AddInt(1)
	output.AddInt(2)
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/mr-tmp-%d.%s",
		dir, 2, consts.InvertedBucket)
	inverted := NewInvertedDB(outName)
	oByte, _ := output.MarshalBinary()
	err := inverted.StoragePostings(query, oByte)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetInvertedInfo(t *testing.T) {
	query := "小岛"
	for i := 0; i < 5; i++ {
		outName := fmt.Sprintf("/Users/mac/GolandProjects/Go-SearchEngine/app/index_platform/woker/mr-tmp-%d.inverted", i)
		inverted := NewInvertedDB(outName)
		oByte, err := inverted.GetInverted([]byte(query))
		if err != nil {
			fmt.Println(err)
		}
		output := roaring.New()
		output.UnmarshalBinary(oByte)
		fmt.Println(output)
	}
}
