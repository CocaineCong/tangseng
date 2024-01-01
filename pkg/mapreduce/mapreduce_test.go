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

package mapreduce_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/RoaringBitmap/roaring"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	analyzer2 "github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/input_data"
	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/ranking"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/pkg/util/stringutils"
	"github.com/CocaineCong/tangseng/types"
)

func TestMain(m *testing.M) {
	analyzer.InitSeg()
	m.Run()
}

func TestMapReduce(t *testing.T) {
	invertedIndex := cmap.New[*roaring.Bitmap]()
	filePaths := []string{"/Users/mac/GolandProjects/Go-SearchEngine/app/mapreduce/input_data/other_input_data/movies_data.csv"}
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		for _, path := range filePaths {
			content, _ := os.ReadFile(path)
			source <- content
		}
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
		res := make([]*types.KeyValue, 0, 1e3)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines[1:] {
			docStruct, _ := doc2Struct(line)
			if docStruct.DocId == 0 {
				continue
			}
			tokens, _ := analyzer2.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
			for _, v := range tokens {
				res = append(res, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
			}
		}
		writer.Write(res)
	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		for values := range pipe {
			for _, v := range values {
				if value, ok := invertedIndex.Get(v.Key); ok {
					value.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, value)
				} else {
					docIds := roaring.NewBitmap()
					docIds.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, docIds)
				}
			}
		}
	})
	keys := invertedIndex.Keys()
	for _, v := range keys {
		val, _ := invertedIndex.Get(v)
		fmt.Println(v, val)
	}
}

func Map(filename string, contents string) (res []*types.KeyValue) {
	res = make([]*types.KeyValue, 0)
	lines := strings.Split(contents, "\r\n")
	for _, line := range lines[1:] {
		docStruct, _ := doc2Struct(line)
		if docStruct.DocId == 0 {
			continue
		}

		tokens, err := analyzer2.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
		if err != nil {
			logs.LogrusObj.Errorf("Map-GseCutForBuildIndex :%+v", err)
			continue
		}
		msgTokens := make([]string, 0)
		for _, v := range tokens {
			res = append(res, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
			msgTokens = append(msgTokens, v.Token)
		}

		// 构建前缀树
		go func(msgTokens []string) {
			err = input_data.DocTrie2Kfk(msgTokens)
			if err != nil {
				return
			}
			defer func() {
				if err := recover(); err != nil {
					logs.LogrusObj.Errorf("input_data.DocTrie2Kfk 消费出现错误 ：%+v", err)
				}
			}()
		}(msgTokens)

		// 建立正排索引
		go func(docStruct *types.Document) {
			err = input_data.DocData2Kfk(docStruct)
			defer func() {
				if err := recover(); err != nil {
					logs.LogrusObj.Errorf("input_data-DocData2Kfk-消费出现错误 :%+v", err)
				}
			}()
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

func TestMapreduceTFIDF(t *testing.T) {
	var searchItem []*types.SearchItem
	for i := 0; i < 3; i++ {
		searchItem = append(searchItem, &types.SearchItem{
			Content: "test" + cast.ToString(i),
			Score:   2 * float64(i),
		})
	}
	words := []string{"test1", "test1"}
	contents := []string{"test", "test"}
	_, err := mapreduce.MapReduce(func(source chan<- *types.SearchItem) {
		for i := range searchItem {
			source <- searchItem[i]
		}
	}, func(item *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		if item != nil {
			for _, word := range words {
				item.Score = ranking.CalculateTFIDF(word, item.Content, contents)
			}
		}
		writer.Write(item)
	}, func(pipe <-chan *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		for values := range pipe {
			if values != nil {
				values.Score += values.Score
			}
			fmt.Println(values)
		}
	})
	fmt.Println(err)
}
