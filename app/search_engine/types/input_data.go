package types

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
