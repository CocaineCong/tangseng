package index

import (
	"fmt"
	"time"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/engine"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/recall"
)

// Recall 召回
type Recall struct {
	*recall.Recall
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
	NewRecallServ(meta).Search(query)
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
