package relevant

import (
	"sort"
	"strings"
	"sync"
)

// 计算词项频率（Term Frequency）
func calculateTF(term string, document string) float64 {
	terms := strings.Fields(document) // 为了方便使用空格分割文档为单词列表
	count := 0
	for _, t := range terms {
		if t == term {
			count++ // 统计词项出现的次数
		}
	}
	return float64(count) / float64(len(terms)) // 返回词项频率（出现次数除以总单词数）
}

// 计算逆文档频率（Inverse Document Frequency）
func calculateIDF(term string, documents []string) float64 {
	docWithTerm := 0
	for _, doc := range documents {
		if strings.Contains(doc, term) {
			docWithTerm++ // 包含term这个词的文档数
		}
	}
	return float64(len(documents)) / float64(docWithTerm)
}

// 计算 BM25
func calculateBM25(term string, document string, documents []string, avgdl, k1 float64, b float64) float64 {
	tf := calculateTF(term, document)                // 计算词项频率
	idf := calculateIDF(term, documents)             // 计算逆文档频率
	doclen := float64(len(strings.Fields(document))) // 文档长度（单词数）

	corpusSize := float64(len(documents)) // 文档集合大小
	df := 0.0
	for _, doc := range documents {
		if strings.Contains(doc, term) {
			df++ // 统计包含给定词项的文档数
		}
	}
	//log := math.Log((corpusSize - df + 0.5) / (df + 0.5))
	bm25 := (idf * tf * (k1 + 1)) / (tf + k1*(1-b+b*(doclen/avgdl))) * ((corpusSize - df + 0.5) / (df + 0.5)) // 计算 BM25 分数
	return bm25
}

// 计算平均文档长度
func calculateAverageDocumentLength(documents []string) float64 {
	totalLength := 0
	for _, doc := range documents {
		totalLength += len(strings.Fields(doc)) // 统计单词数
	}
	return float64(totalLength) / float64(len(documents)) // 返回平均文档长度
}

// 获取分数最高的前 k 个文档的索引
func getTopKDocumentIndices(scores []float64, k int) []int {
	docIndices := make([]int, len(scores))
	for i := range docIndices {
		docIndices[i] = i // 初始化索引数组
	}

	// 按分数降序排序索引
	sort.SliceStable(docIndices, func(i, j int) bool {
		return scores[docIndices[i]] > scores[docIndices[j]]
	})
	//边界检查，以确保k不超过分数数组长度
	if k > len(scores) {
		k = len(scores)
	}

	return docIndices[:k] // 返回前 k 个文档的索引
}

// 并行计算所有文档的BM25分数
func CalculateBM25Scores(query string, documents []string, avgdl float64, k1 float64, b float64) []float64 {
	scores := make([]float64, len(documents))
	var wg sync.WaitGroup
	wg.Add(len(documents))

	ch := make(chan int)
	go func() {
		for i := 0; i < len(documents); i++ {
			idx := <-ch
			scores[idx] = calculateBM25(query, documents[idx], documents, avgdl, k1, b)
			wg.Done()
		}
	}()

	for i := 0; i < len(documents); i++ {
		ch <- i
	}
	wg.Wait()
	close(ch)

	return scores
}
