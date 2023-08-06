package relevant

import (
	"math"
	"sync"

	"github.com/xtgo/set"
)

// ScoreFn is any function that returns a score of the document.
type ScoreFn func(tf *TFIDF, doc Document) []float64

// TFIDF is a structure holding the relevant state information about TF/IDF
type TFIDF struct {
	// Term Frequency
	TF map[int]float64
	// Inverse Document Frequency
	IDF map[int]float64
	// Docs is the count of documents
	Docs int
	// Len is the total length of docs
	Len int
	sync.Mutex
}

// Document is a representation of a document.
type Document interface {
	IDs() []int
}

// New creates a new TFIDF structure
func New() *TFIDF {
	return &TFIDF{
		TF:  make(map[int]float64),
		IDF: make(map[int]float64),
	}
}

// Add adds a document to the state
func (tf *TFIDF) Add(doc Document) {
	ints := BOW(doc)
	tf.Lock()
	for _, w := range ints {
		tf.TF[w]++
	}

	tf.Docs++
	tf.Len += len(ints) // yes we are adding only unique words

	tf.Unlock()
}

// CalculateIDF calculates the inverse document frequency
func (tf *TFIDF) CalculateIDF() {
	docs := float64(tf.Docs)
	tf.Lock()
	for t, f := range tf.TF {
		tf.IDF[t] = math.Log(docs / f)
	}
	tf.Unlock()
}

// TF calculates the term frequencies of term. This is useful for scoring functions.
// It does not make it a unique bag of words.
func TF(doc Document) []float64 {
	ids := doc.IDs()
	retVal := make([]float64, len(ids))
	TF := make(map[int]float64)
	for _, id := range ids {
		TF[id]++
	}

	for i, id := range ids {
		retVal[i] = TF[id]
	}
	return retVal
}

// BOW turns a document into a bag of words. The words of the document will have been deduplicated. A unique list of word IDs is then returned.
func BOW(doc Document) []int {
	ids := doc.IDs()
	retVal := make([]int, len(ids))
	copy(retVal, ids)
	retVal = set.Ints(retVal)
	return retVal
}

// Score calculates the TFIDF score (TF * IDF) for the document without adding the document to the tracked document count.
//
// This function is only useful for a handful of cases. It's recommended you write your own scoring functions.
func (tf *TFIDF) Score(doc Document) []float64 {
	ids := doc.IDs()
	retVal := TF(doc)

	l := float64(len(ids))
	for i, freq := range retVal {
		retVal[i] = (freq / l) * tf.IDF[ids[i]]
	}
	return retVal
}
