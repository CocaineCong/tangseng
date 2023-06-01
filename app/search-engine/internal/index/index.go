package index

import (
	"fmt"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/engine"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/segment"
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
)

type Index struct {
	*engine.Engine
	*engine.Meta
}

// AddDocument 添加文档
func (in *Index) AddDocument(doc *storage.Document) error {
	if doc == nil || doc.DocId <= 0 || doc.Title == "" {
		return fmt.Errorf("doc err: doc || doc_id || title")
	}
	err := in.AddDoc(doc)
	if err != nil {
		return fmt.Errorf("forward doc add err:%v", err)
	}
	err = in.Text2PostingsLists(doc.Title, doc.DocId)
	if err != nil {
		return fmt.Errorf("Text2PostingsLists:%v", err)
	}
	return nil
}

// Close --
func (in *Index) Close() {
	in.Engine.Close()
}

// NewIndexEngine init
func NewIndexEngine(meta *engine.Meta) (*Index, error) {
	e := engine.NewEngine(meta, segment.IndexMode)
	return &Index{
		Engine: e,
		Meta:   meta,
	}, nil
}
