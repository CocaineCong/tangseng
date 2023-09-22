package input_data

import (
	"fmt"

	"github.com/CocaineCong/tangseng/pkg/fileutils"
)

const wukongInputDataPath = "./wukong_input_data"

func WukongData2MapReduce() {
	res := fileutils.GetFiles(wukongInputDataPath)
	fmt.Println(res)
}
