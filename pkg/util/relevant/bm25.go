package relevant

import (
	"sort"

	"github.com/xtgo/set"
)

// DocScore is a tuple of the document ID and a score
type DocScore struct {
	DocId int
	Score float64
}

// DocScores is a list of DocScore
type DocScores []DocScore

func (ds DocScores) Len() int           { return len(ds) }
func (ds DocScores) Less(i, j int) bool { return ds[i].Score < ds[j].Score }
func (ds DocScores) Swap(i, j int) {
	ds[i].Score, ds[j].Score = ds[j].Score, ds[i].Score
	ds[i].DocId, ds[j].DocId = ds[j].DocId, ds[i].DocId
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
