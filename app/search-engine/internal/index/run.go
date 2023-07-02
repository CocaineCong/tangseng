package index

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/engine"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
	"github.com/CocaineCong/Go-SearchEngine/config"
	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
)

func Run(meta *engine.Meta) {
	index, err := NewIndexEngine(meta)
	if err != nil {
		panic(err)
	}
	defer index.Close()

	addDoc(index)
	log.LogrusObj.Infof("index run end")
}

func addDoc(in *Index) {
	// TODO: 后续配置文件改成多选择的
	docList := readFiles([]string{config.Conf.SeConfig.SourceWuKoFile})
	go in.Scheduler.Merge()
	for i, item := range docList[1:] {
		doc, err := doc2Struct(fmt.Sprintf("%d,%s", i, item))
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

	d := strings.Split(docStr, ",")

	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}

	doc := &storage.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[2],
		Body:  d[1],
	}

	return doc, nil
}
