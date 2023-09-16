package input_data

import (
	"fmt"
)

const wukongInputDataPath = "./wukong_input_data"

func WukongData2MapReduce() {
	res := GetFiles(wukongInputDataPath)
	fmt.Println(res)
}
