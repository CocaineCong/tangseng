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

package types

import (
	"github.com/RoaringBitmap/roaring"
)

// Document 文档格式
type Document struct {
	DocId int64  `json:"doc_id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Data2Starrocks struct {
	DocId int64   `json:"doc_id"`
	Url   string  `json:"url"`
	Title string  `json:"title"`
	Desc  string  `json:"desc"`
	Score float64 `json:"score"` // 质量分
}

type Task struct {
	Columns    []string `json:"columns"`
	BiTable    string   `json:"bi_table"`
	SourceType int      `json:"source_type"` // 来源 1 爬虫 2 csv导入
}

type DictTireTree struct {
	Value string `json:"value"`
	Score int64  `json:"score"`
}

// InvertedIndexValue 倒排索引
type InvertedIndexValue struct {
	Token         string        `json:"token"`
	PostingsList  *PostingsList `json:"postings_list"`
	DocCount      int64         `json:"doc_count"`
	PositionCount int64         `json:"position_count"` // 查询使用，写入的时候暂时不用
	TermValues    *TermValue    `json:"term_values"`
}

type TermValue struct {
	DocCount int64 `json:"doc_count"`
	Offset   int64 `json:"offset"`
	Size     int64 `json:"size"`
}

// // PostingsList 倒排列表
// type PostingsList struct {
// 	DocId         int64         `json:"doc_id"`
// 	Positions     []int64       `json:"positions"`
// 	PositionCount int64         `json:"position_count"`
// 	Next          *PostingsList `json:"next"`
// }

type PostingsList struct {
	Term      string          `json:"term"`
	Position  []int64         `json:"position"`   // 位置。为了标红
	TermCount int64           `json:"term_count"` // 个数，为了排序计算，这个词在文档中越多就可能越重要
	DocIds    *roaring.Bitmap `json:"doc_ids"`
}
