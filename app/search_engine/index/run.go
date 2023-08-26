package index

import (
	"fmt"
	"time"

	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	inputData "github.com/CocaineCong/tangseng/app/search_engine/inputdata"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func RunningIndex() {
	meta, err := engine.ParseMeta()
	if err != nil {
		fmt.Println("ParseMeta err", err)
		return
	}
	fmt.Printf("meta: %v", meta)
	// 定时同步meta数据
	ticker := time.NewTicker(time.Minute * 15)
	go meta.SyncByTicker(ticker)
	Run(meta)
	func() {
		// 最后同步元数据至文件
		fmt.Println("close")
		err = meta.SyncMeta()
		if err != nil {
			log.LogrusObj.Errorln("SyncMeta", err)
			return
		}
		fmt.Println("close")
		ticker.Stop()
		fmt.Println("close")
	}()
}

func Run(meta *engine.Meta) {
	index := NewIndexEngine(meta)
	defer index.Close()

	inputData.AddDoc(index)
	log.LogrusObj.Infof("index run end")
}
