package inputData

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search_engine/index"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/stringutils"
)

// AddDoc 读取配置文件，进行doc文件转成struct
func AddDoc(in *index.Index) {
	// TODO: 后续配置文件改成多选择的
	docList := readFiles([]string{config.Conf.SeConfig.SourceWuKoFile})
	go in.Scheduler.Merge()
	wg := new(sync.WaitGroup)
	for _, item := range docList[1:] {
		wg.Add(1)
		go func(item string) {
			doc, err := doc2Struct(item)
			if err != nil {
				log.LogrusObj.Errorf("index addDoc doc2Struct: %v", err)
			}
			err = in.AddDocument(doc)
			if err != nil {
				log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
			}
			wg.Done()
		}(item)
	}
	wg.Wait()
	// 读取结束 写入磁盘
	err := in.FlushInvertedIndex(true)
	if err != nil {
		log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
		return
	}
}

// doc2Struct 从csv读取数据 TODO：后续区分一下输入源，如果是爬虫那边的数据，处理不一样
func doc2Struct(docStr string) (*types.Document, error) {
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	if len(d) < 3 {
		return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	}

	doc := &types.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  stringutils.StrConcat([]string{d[2], d[3], d[4]}),
	}

	return doc, nil
}
