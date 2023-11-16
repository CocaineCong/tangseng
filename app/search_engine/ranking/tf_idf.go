package ranking

import (
	"fmt"
	"math"
	"sync"

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

// 计算TF-IDF值
func CalculateTFIDF(word string, doc string, docs []string) float64 {
	tf := calculateTF(word, doc)
	idf := calculateIDF(word, docs)
	return tf * idf
}

func CalculateScoreTFIDF(token string, searchItem []*types.SearchItem) []*types.SearchItem {
	contents := make([]string, len(searchItem))
	for i, item := range searchItem {
		contents[i] = item.Content
	}
	words := analyzer.GlobalSega.Cut(token, true)
	wg := new(sync.WaitGroup)
	for i := range searchItem {
		wg.Add(1)
		go func(i int) {
			for _, word := range words {
				searchItem[i].Score += CalculateTFIDF(word, searchItem[i].Content, contents)
			}
			wg.Done()
		}(i)
	}

	mapreduce.MapReduce(func(source chan<- *types.SearchItem) {
		for i := range searchItem {
			source <- searchItem[i]
		}
	}, func(item *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		for _, word := range words {
			item.Score = CalculateTFIDF(word, item.Content, contents)
		}
		writer.Write(item)
	}, func(pipe <-chan *types.SearchItem, writer mapreduce.Writer[*types.SearchItem], cancel func(err error)) {
		for values := range pipe {
			fmt.Println("values", values)
			values.Score += values.Score
		}
	})

	wg.Wait()

	return searchItem
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
