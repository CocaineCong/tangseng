package relevant

import (
	"fmt"

	"testing"
)

// 基准测试
func BenchmarkBM25T(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BM25T()
	}
}

// 基准测试用例
func BM25T() {
	bodyRecallReason = []string{
		"This is the first document",
		"This document is the second document",
		"And This is the third one",
		"Is this the first document?",
		"This This This",
	}

	corpus, _ := MakeCorpus(bodyRecallReason)
	docs := MakeDocuments(bodyRecallReason, corpus)
	tf := New()

	for _, doc := range docs {
		tf.Add(doc)
	}
	tf.CalculateIDF()
	token := Doc{corpus["This"]}
	tokenScores := BM25(tf, token, docs, 1.5, 0.75)

	for _, d := range tokenScores[:3] {
		if d.Score == 0.0 {
			continue
		}
	}
}

func TestBm25direct(t *testing.T) {
	// 示例文档集合
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
	avgdl := calculateAverageDocumentLength(documents) // 计算平均文档长度

	// 执行查询并获取结果
	scores := CalculateBM25Scores(query, documents, avgdl, k1, b)

	// 按分数排序并返回前 k 个文档的索引
	docIndices := getTopKDocumentIndices(scores, 5)

	// 输出查询结果
	for _, index := range docIndices {
		fmt.Println("Document: ", documents[index], "分数：", scores[index]) // 输出查询结果
	}
}

// 基准测试
func BenchmarkBm25direct(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bm25direct()

	}
}

// 基准测试用例
func Bm25direct() {
	// 示例文档集合
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
	avgdl := calculateAverageDocumentLength(documents) // 计算平均文档长度

	// 执行查询并获取结果
	_ = CalculateBM25Scores(query, documents, avgdl, k1, b)

}

// bm25direct:   117639	     10520 ns/op
// bm25      :     4615	    278839 ns/op
// --- PASS: TestBM (3.08s)
func TestBM(t *testing.T) {
	// 运行基准测试，并产生报告
	bm25directRootResult := testing.Benchmark(BenchmarkBm25direct)
	bm25TRootResult := testing.Benchmark(BenchmarkBM25T)
	fmt.Printf("bm25direct: %s\n", bm25directRootResult.String())
	fmt.Printf("bm25      : %s\n", bm25TRootResult.String())
}
