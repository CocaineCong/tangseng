package query

import (
	"fmt"
	"testing"

	"github.com/go-ego/gse"
)

func TestQuery(t *testing.T) {
	newGse, _ := gse.New()
	// content := "疑电影《擒凶记1894》由香港电影人陈德森和陈锡康合作，讲述19世纪末的香港无头悬案，从而引发出一桩惊天阴谋，来自英国一丝不苟的年轻法医，与世袭嘉业的粗鄙侩子手联袂搭档，在追查无头命案的背后，有着令人不寒而栗的冷血真相。该片目前演员阵容还未确认曝光，据悉有望重金聘请国外巨星投身剧组。"
	content := "机械"
	content = ignoredChar(content)
	segments := newGse.Pos(content)
	for i := range segments {
		fmt.Println(segments[i])
	}

}
