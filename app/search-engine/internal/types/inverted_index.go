package types

import (
	"github.com/CocaineCong/tangseng/app/search-engine/internal/storage"
)

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      int64
	PositionCount int64 // 查询使用，写入的时候暂时不用
	TermValues    *storage.TermValue
}

// PostingsList 倒排列表
type PostingsList struct {
	DocId         int64
	Positions     []int64
	PositionCount int64
	Next          *PostingsList
}
