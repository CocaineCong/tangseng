package inputData

import (
	"strings"
	"sync"

	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/search_engine/index"
	"github.com/CocaineCong/tangseng/app/search_engine/types"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// AddDoc 读取配置文件，进行doc文件转成struct
func AddDoc(in *index.Index) {
	// TODO: 后续配置文件改成多选择的
	docList := readFiles([]string{config.Conf.SeConfig.SourceWuKoFile})
	go in.Scheduler.Merge()
	wg := new(sync.WaitGroup)
	for _, item := range docList[1:] {
		wg.Add(1)
		doc, err := doc2Struct(item)
		if err != nil {
			log.LogrusObj.Errorf("index addDoc doc2Struct: %v", err)
			continue
		}
		err = in.AddDocument(doc)
		if err != nil {
			log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
			continue
		}
		wg.Done()
	}
	wg.Wait()
	// 读取结束 写入磁盘
	err := in.Flush(true)
	if err != nil {
		log.LogrusObj.Errorf("index addDoc AddDocument: %v", err)
		return
	}
}

func doc2Struct(docStr string) (*types.Document, error) {
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	// if len(d) < 3 { // TODO: 后续记得开放
	// 	return nil, fmt.Errorf("doc2Struct err: %v", "docStr is not right")
	// }

	doc := &types.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  d[1],
	}

	return doc, nil
}
