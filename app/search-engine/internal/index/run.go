package index

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/engine"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func IndexRunning() {
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
		meta.SyncMeta()
		fmt.Println("close")
		ticker.Stop()
		fmt.Println("close")
	}()
}

func Run(meta *engine.Meta) {
	index := NewIndexEngine(meta)
	defer index.Close()

	addDoc(index)
	log.LogrusObj.Infof("index run end")
}

func addDoc(in *Index) {
	// TODO: 后续配置文件改成多选择的
	docList := readFiles([]string{config.Conf.SeConfig.SourceWuKoFile})
	go in.Scheduler.Merge()
	for _, item := range docList[1:40] {
		doc, err := doc2Struct(item)
		if err != nil {
			log.LogrusObj.Errorf("index addDoc doc2Struct: %v", err)
		}
		err = in.AddDocument(doc)
		if err != nil {
			log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
		}
	}
	// 读取结束 写入磁盘
	err := in.Flush(true)
	if err != nil {
		log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
	}
}

func doc2Struct(docStr string) (*storage.Document, error) {
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}

	doc := &storage.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  d[1],
	}
	fmt.Println("doc", doc.DocId, doc.Body)

	return doc, nil
}
