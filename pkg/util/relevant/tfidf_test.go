package relevant

import (
	"fmt"
	"testing"
)

func TestTFIDF(t *testing.T) {
	corpus, _ := MakeCorpus(bodyRecallReason)
	docs := MakeDocuments(bodyRecallReason, corpus)
	tf := New()

	for _, doc := range docs {
		tf.Add(doc)
	}

	tf.CalculateIDF()
	for _, doc := range docs {
		score := tf.Score(doc)
		fmt.Println(score)
	}
}
