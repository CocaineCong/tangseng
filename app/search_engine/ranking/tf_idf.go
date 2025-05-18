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

package ranking

import (
	"math"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/types"
)

// 计算TF值
func calculateTF(word string, doc string) float64 {
	words := analyzer.GlobalSega.Cut(doc, true)
	count := 0
	for _, w := range words {
		if w == word {
			count++
		}
	}
	return float64(count) / float64(len(words))
}

// 计算IDF值
func calculateIDF(word string, docs []string) float64 {
	count := 0
	for _, doc := range docs {
		words := analyzer.GlobalSega.Cut(doc, true)
		for _, w := range words {
			if w == word {
				count++
				break
			}
		}
	}
	return math.Log(float64(len(docs)) / float64(count+1))
}

// CalculateTFIDF 计算TF-IDF值
func CalculateTFIDF(word string, doc string, docs []string) float64 {
	tf := calculateTF(word, doc)
	idf := calculateIDF(word, docs)
	return tf * idf
}

func CalculateScoreTFIDF(token string, searchItem []*types.SearchItem) (resp []*types.SearchItem, err error) {
	contents := make([]string, len(searchItem))
	resp = make([]*types.SearchItem, len(searchItem))
	for i, item := range searchItem {
		contents[i] = item.Content
	}
	words := analyzer.GlobalSega.Cut(token, true)
	index := 0
	_, err = mapreduce.MapReduce(func(source chan<- *types.SearchItem) {
		for i := range searchItem {
			source <- searchItem[i]
		}
	}, func(item *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		if item != nil {
			for _, word := range words {
				item.Score = CalculateTFIDF(word, item.Content, contents)
			}
		}
		writer.Write(item)
	}, func(pipe <-chan *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		for values := range pipe {
			if values != nil {
				values.Score += values.Score
			}
			resp[index] = values
			index++
		}
	})

	return
}

// func CalculateScoreTFIDF(token string, searchItem []*types.SearchItem) (resp []*types.SearchItem) {
// 	contents := make([]string, 0)
// 	for i := range searchItem {
// 		contents = append(contents, searchItem[i].Content)
// 	}
// 	words := analyzer.GlobalSega.Cut(token, true)
// 	for i := range searchItem {
// 		for _, word := range words {
// 			searchItem[i].Score += calculateTFIDF(word, searchItem[i].Content, contents)
// 		}
// 	}

// 	resp = searchItem
// 	return
// }
