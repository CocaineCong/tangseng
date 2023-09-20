package types

// SearchItem 查询结果
type SearchItem struct {
	DocId        int64   `json:"doc_id"`
	Content      string  `json:"content"`
	Title        string  `json:"title"`
	Score        float64 `json:"score"`         // 这个词对于这篇文章的评分，也就是这个词到底重不重要
	DocCount     int64   `json:"doc_count"`     // 这个词在文中出现了多少次
	ContentScore float64 `json:"content_score"` // 这篇文章的评分
}

type SearchItemList []*SearchItem

func (ds SearchItemList) Len() int           { return len(ds) }
func (ds SearchItemList) Less(i, j int) bool { return ds[i].Score < ds[j].Score }
func (ds SearchItemList) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

// 用于实现排序的map
type queryTokenHash struct {
	token         string
	invertedIndex *InvertedIndexValue
	fetchPostings *PostingsList
}

// token游标 标识当前位置
type searchCursor struct {
	doc     *PostingsList // 文档编号的序列
	current *PostingsList // 当前文档编号
}

// 短语游标
type phraseCursor struct {
	positions []int64 // 位置信息
	base      int64   // 词元在查询中的位置
	current   *int64  // 当前的位置信息
	index     int     // 当前位置index
}
