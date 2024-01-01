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

package relevant

import (
	"sort"

	"github.com/xtgo/set"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
)

// DocScore is a tuple of the document ID and a score
type DocScore struct {
	ID    int
	Score float64
}

// DocScores is a list of DocScore
type DocScores []DocScore

func (ds DocScores) Len() int           { return len(ds) }
func (ds DocScores) Less(i, j int) bool { return ds[i].Score < ds[j].Score }
func (ds DocScores) Swap(i, j int) {
	ds[i].Score, ds[j].Score = ds[j].Score, ds[i].Score
	ds[i].ID, ds[j].ID = ds[j].ID, ds[i].ID
}

// BM25 is the scoring function.
//
// k1 should be between 1.2 and 2.
// b should be around 0.75
func BM25(tf *TFIDF, query Document, docs []Document, k1, b float64) DocScores {
	q := BOW(query)
	w := make([]int, len(q))
	copy(w, q)
	avgLen := float64(tf.Len) / float64(tf.Docs)

	scores := make([]float64, 0, len(docs))
	for _, doc := range docs {
		// TF := tfidf.TF(docs)
		d := BOW(doc)
		w = append(w, d...)
		size := set.Inter(sort.IntSlice(w), len(q))
		n := w[:size]

		score := make([]float64, 0, len(n))
		docLen := float64(len(d))
		for _, id := range n {
			num := tf.TF[id] * (k1 + 1)
			denom := tf.TF[id] + k1*(1-b+b*docLen/avgLen)
			idf := tf.IDF[id]
			score = append(score, idf*num/denom)
		}
		scores = append(scores, sum(score))

		// reset working vector
		copy(w, q)
		w = w[:len(q)]
	}
	var retVal DocScores
	for i := range docs {
		retVal = append(retVal, DocScore{i, scores[i]})
	}
	return retVal
}

func sum(a []float64) float64 {
	var retVal float64
	for _, f := range a {
		retVal += f
	}
	return retVal
}

type Doc []int

func (d Doc) IDs() []int { return d }

func MakeCorpus(a []string) (map[string]int, []string) {
	retVal := make(map[string]int)
	invRetVal := make([]string, 0)
	var id int
	for _, s := range a {
		tokens, _ := analyzer.GseCutForRecall(s)
		for _, f := range tokens {
			if _, ok := retVal[f]; !ok {
				retVal[f] = id
				invRetVal = append(invRetVal, f)
				id++
			}
		}
	}
	return retVal, invRetVal
}

func MakeDocuments(a []string, c map[string]int) []Document {
	retVal := make([]Document, 0, len(a))
	for _, s := range a {
		var ts []int
		tokens, _ := analyzer.GseCutForRecall(s)
		for _, f := range tokens {
			id := c[f]
			ts = append(ts, id)
		}
		retVal = append(retVal, Doc(ts))
	}
	return retVal
}
