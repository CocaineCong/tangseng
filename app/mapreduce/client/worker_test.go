package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/mapreduce/analyzer"
	"github.com/CocaineCong/tangseng/app/mapreduce/rpc"
	"github.com/CocaineCong/tangseng/app/mapreduce/service/input_data_mr"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	analyzer.InitSeg()
	rpc.Init()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestWorker(t *testing.T) {
	ctx := context.Background()
	Worker(ctx, input_data_mr.Map, input_data_mr.Reduce)
}
