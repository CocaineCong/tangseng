package types

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string        `json:"token"`
	PostingsList  *PostingsList `json:"postings_list"`
	DocCount      int64         `json:"doc_count"`
	PositionCount int64         `json:"position_count"` // 查询使用，写入的时候暂时不用
	TermValues    *TermValue    `json:"term_values"`
}

type TermValue struct {
	DocCount int64 `json:"doc_count"`
	Offset   int64 `json:"offset"`
	Size     int64 `json:"size"`
}

// PostingsList 倒排列表
type PostingsList struct {
	DocId         int64         `json:"doc_id"`
	Positions     []int64       `json:"positions"`
	PositionCount int64         `json:"position_count"`
	Next          *PostingsList `json:"next"`
}
