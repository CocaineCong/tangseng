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
	"sort"

	"github.com/CocaineCong/tangseng/pkg/util/relevant"
	"github.com/CocaineCong/tangseng/types"
)

// CalculateScoreBm25 计算相关性
func CalculateScoreBm25(token string, searchItem []*types.SearchItem) (resp []*types.SearchItem) {
	contents := make([]string, 0)
	for i := range searchItem {
		contents = append(contents, searchItem[i].Content)
	}
	corpus, _ := relevant.MakeCorpus(contents)
	docs := relevant.MakeDocuments(contents, corpus)
	tf := relevant.New()
	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()
	tokenRecall := relevant.Doc{corpus[token]}
	bm25Scores := relevant.BM25(tf, tokenRecall, docs, 1.5, 0.75)
	sort.Sort(sort.Reverse(bm25Scores))

	for i := range bm25Scores {
		if bm25Scores[i].Score == 0.0 {
			continue
		}
		searchItem[bm25Scores[i].ID].Score = bm25Scores[i].Score
	}
	sort.Slice(searchItem, func(i, j int) bool {
		return searchItem[i].Score > searchItem[j].Score
	})

	resp = searchItem

	return
}
