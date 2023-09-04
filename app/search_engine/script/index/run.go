package index

import (
	"context"
	"time"

	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func RunningIndex() {
	ctx := context.Background()
	meta, err := engine.ParseMeta()
	if err != nil {
		log.LogrusObj.Errorln("ParseMeta err", err)
		return
	}
	log.LogrusObj.Errorf("meta: %v", meta)
	// 定时同步meta数据
	ticker := time.NewTicker(time.Minute * 15)
	go meta.SyncByTicker(ticker)
	Run(ctx, meta)
	SyncIndex2Meta(meta, ticker)
}

func Run(ctx context.Context, meta *engine.Meta) {
	in := NewIndexEngine(meta)
	defer in.Close()

	AddDoc(ctx, in)
	log.LogrusObj.Infof("index run end")
}

func SyncIndex2Meta(meta *engine.Meta, ticker *time.Ticker) {
	// 最后同步元数据至文件
	log.LogrusObj.Infof("close")
	err := meta.SyncMeta()
	if err != nil {
		log.LogrusObj.Errorln("SyncIndex2Meta-SyncMeta", err)
		return
	}
	log.LogrusObj.Infof("close")
	ticker.Stop()
	log.LogrusObj.Infof("close")
}
