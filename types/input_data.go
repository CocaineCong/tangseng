package types

type InputDataList struct {
	DocId int64   `json:"doc_id"`
	Title string  `json:"title"`
	Url   string  `json:"url"`
	Body  string  `json:"body"`
	Score float64 `json:"score"`
}
