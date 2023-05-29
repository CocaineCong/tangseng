package segment

import (
	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/storage"
	"github.com/CocaineCong/Go-SearchEngine/config"
)

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string
	PostingsList  *PostingsList
	DocCount      int64
	PositionCount int64 // 查询使用，写入的时候暂时不用
	TermValues    *storage.TermValue
}

// InvertedIndexHash 倒排hash
type InvertedIndexHash map[string]*InvertedIndexValue

// CreateNewInvertedIndex 创建倒排索引
func CreateNewInvertedIndex(token string, docCount int64) *InvertedIndexValue {
	p := new(InvertedIndexValue)
	p.DocCount = docCount
	p.Token = token
	p.PositionCount = 0
	p.PostingsList = new(PostingsList)
	return p
}

// SegmentDbInit 读取对应segment文件下的db
func SegmentDbInit(segId SegId, conf *config.Config) {

}
