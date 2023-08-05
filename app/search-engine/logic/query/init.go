package query

import (
	"github.com/go-ego/gse"
)

var (
	GobalSeg gse.Segmenter
)

func InitSeg() {
	newGse, _ := gse.New()
	GobalSeg = newGse
}
