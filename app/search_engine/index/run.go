package index

import (
	"fmt"
	"time"

	"github.com/CocaineCong/tangseng/app/search_engine/engine"
	inputData "github.com/CocaineCong/tangseng/app/search_engine/inputdata"
	"github.com/CocaineCong/tangseng/config"
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
	in := NewIndexEngine(meta)
	defer in.Close()

	AddDoc(in)
	log.LogrusObj.Infof("index run end")
}

// AddDoc 读取配置文件，进行doc文件转成struct
func AddDoc(in *Index) {
	// TODO: 后续配置文件改成多选择的
	docList := inputData.ReadFiles([]string{config.Conf.SeConfig.SourceWuKoFile})
	go in.Scheduler.Merge()
	// wg := new(sync.WaitGroup)
	for _, item := range docList[1:20] {
		// wg.Add(1)
		// go func(item string) {
		doc, err := inputData.Doc2Struct(item)
		if err != nil {
			log.LogrusObj.Errorf("index addDoc doc2Struct: %v", err)
		}

		err = in.AddDocument(doc)
		// }(item)
	}
	// wg.Wait()
	// 读取结束 写入磁盘
	err := in.FlushDict(true)
	if err != nil {
		log.LogrusObj.Errorf("AddDoc-FlushDict: %v", err)
		return
	}

	err = in.FlushInvertedIndex(true)
	if err != nil {
		log.LogrusObj.Errorf("AddDoc-FlushInvertedIndex: %v", err)
		return
	}

}
