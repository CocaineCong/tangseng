package analyzer

import (
	"github.com/go-ego/gse"
)

// 分词器相关

var (
	GlobalSega gse.Segmenter
)

// InitSeg 分词器初始化
func InitSeg() {
	newGse, _ := gse.New()
	GlobalSega = newGse
}
