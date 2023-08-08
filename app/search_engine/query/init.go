package query

import (
	"github.com/go-ego/gse"
)

var (
	GobalSeg gse.Segmenter
)

// InitSeg 分词器初始化
func InitSeg() {
	newGse, _ := gse.New()
	GobalSeg = newGse
}
