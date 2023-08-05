package test

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/index"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/query"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	query.InitSeg()
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestRecall(t *testing.T) {
	q := "电影"
	index.SearchRecall(q)
}
