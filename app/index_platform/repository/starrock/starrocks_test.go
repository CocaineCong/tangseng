package starrock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

// mysql -h 127.0.0.1 -P9030 -uroot
func TestDirectUpload_StreamUpload(t *testing.T) {
	ctx := context.Background()
	du := NewDirectUpload(ctx, &types.Task{
		Columns:    []string{"doc_id", "url", "title", "desc", "score"},
		BiTable:    "test_upload",
		SourceType: 0,
	})
	fmt.Println(config.Conf.StarRocks)
	du.Push(&types.Data2Starrocks{
		DocId: 1,
		Url:   "https://localhost:8083",
		Title: "这是一个测试文件",
		Desc:  "进行测试作用",
		Score: 1220.120,
	})
	time.Sleep(10 * time.Second)
}
