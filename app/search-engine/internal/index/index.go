package index

import (
	"fmt"

	"github.com/CocaineCong/tangseng/app/search-engine/internal/engine"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/segment"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type Index struct {
	*engine.Engine
	*engine.Meta
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
func NewIndexEngine(meta *engine.Meta) *Index {
	return &Index{
		Engine: engine.NewEngine(meta, segment.IndexMode),
		Meta:   meta,
	}
}
