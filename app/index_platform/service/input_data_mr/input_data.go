package input_data_mr

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/input_data"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/util/stringutils"
	"github.com/CocaineCong/tangseng/types"
)

func Map(filename string, contents string) (res []*types.KeyValue) {
	res = make([]*types.KeyValue, 0)
	lines := strings.Split(contents, "\r\n")
	for _, line := range lines[1:] {
		docStruct, _ := doc2Struct(line)
		if docStruct.DocId == 0 {
			continue
		}

		tokens, err := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
		if err != nil {
			logs.LogrusObj.Errorf("Map-GseCutForBuildIndex :%+v", err)
			continue
		}

		for _, v := range tokens {
			res = append(res, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
			go func(token string) {
				input_data.DocTrie2Kfk(token)
			}(v.Token)
		}

		// 存正排索引库中
		// inputData = &model.InputData{
		// 	DocId:   docStruct.DocId,
		// 	Title:   docStruct.Title,
		// 	Body:    docStruct.Body,
		// 	Url:     "",
		// 	Score:   float64(len(docStruct.Body)),
		// 	Source:  consts.DataSourceCSV,
		// 	IsIndex: true,
		// }
		// 建立正排索引
		go func(docStruct *types.Document) {
			input_data.DocData2Kfk(docStruct)
		}(docStruct)
	}

	return
}

func Reduce(key string, values []string) *roaring.Bitmap {
	docIds := roaring.New()
	for _, v := range values {
		docIds.AddInt(cast.ToInt(v))
	}
	return docIds
}

func doc2Struct(docStr string) (doc *types.Document, err error) {
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	something2Str := make([]string, 0)

	for i := 2; i < 5; i++ {
		if len(d) > i && d[i] != "" {
			something2Str = append(something2Str, d[i])
		}
	}

	doc = &types.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  stringutils.StrConcat(something2Str),
	}

	return
}
