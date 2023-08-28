package analyzer

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-ego/gse"

	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	InitSeg()
	log.InitLog()
	m.Run()
}

func TestQuery(t *testing.T) {
	newGse, _ := gse.New()
	// content := "疑电影《擒凶记1894》由香港电影人陈德森和陈锡康合作，讲述19世纪末的香港无头悬案，从而引发出一桩惊天阴谋，来自英国一丝不苟的年轻法医，与世袭嘉业的粗鄙侩子手联袂搭档，在追查无头命案的背后，有着令人不寒而栗的冷血真相。该片目前演员阵容还未确认曝光，据悉有望重金聘请国外巨星投身剧组。"
	content := "时代少年团"
	// content := "机械"
	content = ignoredChar(content)
	segments := newGse.CutAll(content)
	// segments := newGse.Pos(content, true)
	// segments := newGse.CutSearch(content, true)
	for i := range segments {
		// fmt.Println(segments[i].Token().Text(), segments[i].Start(), segments[i].End())
		fmt.Println(segments[i])
	}
}

func TestGseCut(t *testing.T) {
	fileName := "../data/movies.csv"
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	docList := strings.Split(string(content), "\n")
	if len(docList) == 0 {
		fmt.Println(err)
	}
	for _, v := range docList[1:] {
		// tokens, _ := GseCut(v)
		tm := strings.Split(v, ",")
		if len(tm) >= 2 {
			tokens, _ := GseCut(tm[1])
			fmt.Println(tokens)
		}
	}
}
