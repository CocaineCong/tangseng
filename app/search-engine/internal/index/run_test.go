package index

import (
	"fmt"
	"testing"

	"github.com/CocaineCong/Go-SearchEngine/config"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestRunIndex(t *testing.T) {

}
