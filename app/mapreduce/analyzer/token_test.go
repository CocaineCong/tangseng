package analyzer

import (
	"fmt"
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
	// contents := []string{"时代少年团", "冬季卫衣推荐", "小岛秀夫", "我有一个小东西"}
	contents := []string{"动作爱情桀骜不驯的如龙武功高强在一场比赛中被打成重伤被诊今生不能再用功夫女友荆兰为激发他重新振作按照他的原型制造了一个机器人如龙能否恢复武功他和机器人到底哪个功夫高强又能否在机器人中迷失自己穿越华语"}
	for _, content := range contents {
		fmt.Println("*******CutAll")
		content = ignoredChar(content)
		cutAllSegments := newGse.CutAll(content)
		for i := range cutAllSegments {
			fmt.Println(cutAllSegments[i])
		}

		fmt.Println("*******CutSearch")
		searchSegment := newGse.CutSearch(content, true)
		for _, v := range searchSegment {
			fmt.Println(v)
		}

		fmt.Println("*******Segment")
		segmentSegment := newGse.Segment([]byte(content))
		for _, v := range segmentSegment {
			fmt.Println(v.Token().Text())
		}
	}

}
