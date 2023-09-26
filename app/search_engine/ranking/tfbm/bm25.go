package tfbm

import (
	"strings"
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
func CalculateBM25(term string, document string, documents []string, k1 float64, b float64) float64 {
	tf := calculateTF(term, document)                  // 计算词项频率
	idf := calculateIDF(term, documents)               // 计算逆文档频率
	doclen := float64(len(strings.Fields(document)))   // 文档长度（单词数）
	avgdl := calculateAverageDocumentLength(documents) // 计算平均文档长度
	corpusSize := float64(len(documents))              // 文档集合大小
	df := 0.0
	for _, doc := range documents {
		if strings.Contains(doc, term) {
			df++ // 统计包含给定词项的文档数
		}
	}
	//实际应用时,取自然对数，可以使得idf的取值范围更广泛，更好地表征词项的稀有程度。
	//log := math.Log((corpusSize - df + 0.5) / (df + 0.5))
	//bm25 := (idf * tf * (k1 + 1)) / (tf + k1*(1-b+b*(doclen/avgdl))) * log
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
