package storage

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/query"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	query.InitSeg()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestForwardDBRead(t *testing.T) {
	a := config.Conf.SeConfig.StoragePath + "0.forward"
	forward, err := NewForwardDB(a)
	if err != nil {
		fmt.Println("err", err)
	}
	count, err := forward.ForwardCount()
	if err != nil {
		fmt.Println("Err", err)
	}
	fmt.Println(count)
	r, err := forward.GetForward(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(r))
}
