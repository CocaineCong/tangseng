package tfbm

import (
	"testing"
)

func BenchmarkSprint(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bm25()

	}
}

func Bm25() {
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

	}

}

