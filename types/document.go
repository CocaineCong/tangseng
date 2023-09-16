package types

// Document 文档格式
type Document struct {
	DocId int64  `json:"doc_id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Data2Starrocks struct {
	DocId int64   `json:"doc_id"`
	Url   string  `json:"url"`
	Title string  `json:"title"`
	Desc  string  `json:"desc"`
	Score float64 `json:"score"` // 质量分
}

type Task struct {
	Columns    []string `json:"columns"`
	BiTable    string   `json:"bi_table"`
	SourceType int      `json:"source_type"` // 来源 1 爬虫 2 csv导入
}

type DictTireTree struct {
	Value string `json:"value"`
	Score int64  `json:"score"`
}

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
