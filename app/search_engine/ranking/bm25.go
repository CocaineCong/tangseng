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
	resp = make([]*types.SearchItem, 0)
	resp = searchItem

	return
}
