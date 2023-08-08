package index

import (
	"fmt"

	engine2 "github.com/CocaineCong/tangseng/app/search_engine/engine"
	"github.com/CocaineCong/tangseng/app/search_engine/segment"
	"github.com/CocaineCong/tangseng/app/search_engine/storage"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type Index struct {
	*engine2.Engine
	*engine2.Meta
}

// AddDocument 添加文档
func (in *Index) AddDocument(doc *storage.Document) (err error) {
	if doc == nil || doc.DocId <= 0 || doc.Title == "" {
		return fmt.Errorf("doc err: doc || doc_id || title")
	}
	err = in.AddDoc(doc)
	if err != nil {
		logs.LogrusObj.Errorf("forward doc add err:%v", err)
		return
	}
	err = in.Text2PostingsLists(doc.Body, doc.DocId)
	if err != nil {
		logs.LogrusObj.Errorf("Text2PostingsLists:%v", err)
		return
	}

	return
}

// Close --
func (in *Index) Close() {
	in.Engine.Close()
}

// NewIndexEngine init
func NewIndexEngine(meta *engine2.Meta) *Index {
	return &Index{
		Engine: engine2.NewEngine(meta, segment.IndexMode),
		Meta:   meta,
	}
}
