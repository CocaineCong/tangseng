package model

type InputData struct {
	Id      int64 `gorm:"primarykey"`
	DocId   int64 `gorm:"index"`
	Title   string
	Body    string
	Url     string
	Score   float64
	Source  int
	IsIndex bool
}
