package inputData

type InputData struct {
	Key  string      `json:"key"`
	Data interface{} `json:"data"`
}

//func Insert (inData *InputData) (bool, error) {
//	indexSet := index.GetIndexSet() //后续程序应该是常驻内存的，避免io反复操作
//	key := inData.Key
//	if key == "" {
//		return false, errors.New("插入数据的Key不能为空字符串")
//	}
//
//	iData, ok := inData.Data.(map[string]interface{})
//	if !ok {//断言失败应该报错
//		return false, errors.New("插入数据的类型必须为字符串类型的键值对，键值都必须为字符串，不能有复杂层级,收到：" + reflect.TypeOf(inData.Data).String())
//	}
//
//	//整理本次会出现的所有索引
//	for i, _ := range iData {
//		indexSet.Sets = append(indexSet.Sets, i)//带上上次的结果
//	}
//
//	//所有会出现的索引
//	indexSet.Sets = utils.ArrayUnique(indexSet.Sets)
//	indexSet.Save()
//
//	//构建倒排索引
//	for _, revIndexName := range indexSet.Sets {
//		revIndex := index.GetRevIndex(revIndexName)
//		value, ok := iData[revIndexName].(string)
//		if !ok {//断言失败
//			return false, errors.New("插入数据的类型必须为字符串类型的键值对，键值都必须为字符串，不能有复杂层级,收到：" + reflect.TypeOf(iData[revIndexName]).String())
//		}
//
//		if _, ok := revIndex.Data[value];!ok {
//			revIndex.Data[value] = []string{}
//		}
//
//		revIndex.Data[value] = append(revIndex.Data[value], key)
//		revIndex.Data[value] = utils.ArrayUnique(revIndex.Data[value])
//		revIndex.Save()
//	}
//
//	return true, nil
//}
