package tfbm

import (
	"fmt"
	"sort"
	"testing"
)

func TestBm(t *testing.T) {
	//示例文档集合
	documents := []string{
		"This is the first document",
		"This document is the second document",
		"And This is the third one",
		"Is this the first document?",
		"This This This",
	}

	// 查询文本和相关参数
	query := "This"
	k1 := 1.2
	b := 0.75

	// 执行查询
	scores := make([]float64, len(documents))
	for i, doc := range documents {
		scores[i] = CalculateBM25(query, doc, documents, k1, b) // 计算每个文档的 BM25 分数
		fmt.Println("文档:", doc, "---分数 :", scores[i])
	}

	// 按分数排序并返回前 k 个文档的索引
	docIndices := getTopKDocumentIndices(scores, 5)

	// 输出查询结果
	for _, index := range docIndices {
		fmt.Printf("Document: %s\n", documents[index]) // 输出查询结果
	}
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

	return docIndices[:k] // 返回前 k 个文档的索引
}
