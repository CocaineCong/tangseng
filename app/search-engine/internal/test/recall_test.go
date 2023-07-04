package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/engine"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/recall"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// Recall 召回
type Recall struct {
	*recall.Recall
}

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

// NewRecallServ 创建召回服务
func NewRecallServ(meta *engine.Meta) *Recall {
	r := recall.NewRecall(meta)
	return &Recall{r}
}

func SearchRecall(query string) {
	meta, err := engine.ParseMeta()
	if err != nil {
		panic(err)
	}

	// 定时同步meta数据
	ticker := time.NewTicker(time.Second * 10)
	go meta.SyncByTicker(ticker)
	recall := NewRecallServ(meta)
	recall.Search(query)
	// close
	func() {
		// 最后同步元数据至文件
		fmt.Println("close")
		meta.SyncMeta()
		fmt.Println("close")
		ticker.Stop()
		fmt.Println("close")
	}()
}

func TestRecall(t *testing.T) {
	query := "英超比赛"
	SearchRecall(query)
}
