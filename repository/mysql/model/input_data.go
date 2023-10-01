package model

type InputData struct {
	Id      int64  `gorm:"primarykey"`
	DocId   int64  `gorm:"index"`
	Title   string `gorm:"type:longtext"`
	Body    string `gorm:"type:longtext"`
	Url     string
	Score   float64
	Source  int
	IsIndex bool
}
