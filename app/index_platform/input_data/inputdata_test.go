package input_data

import (
	"fmt"
	"testing"
)

func TestInputDataDoc2Struct(t *testing.T) {
	a := "2,[安乐乡]导演利桑德罗·阿隆索导演将打造下一部影片[尤里卡](Eureka，暂译)。据悉该片探讨美国文化问题，故事发生在1870年到2019年期间，涉及地区包括美国、墨西哥以及亚马逊雨林。故事主角是一个经历波折，辗转各地的女性。本片今年7月已在达科他开拍，预计将在2020年上映。"
	r, err := Doc2Struct(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("r", r)
}
