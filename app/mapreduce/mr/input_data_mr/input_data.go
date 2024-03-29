// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package input_data_mr

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
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
		}
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
