package pr

import (
	"math"
	"testing"
)

func TestPageRank(t *testing.T) {
	g := NewGraph()
	g.AddLink("A", "B")
	g.AddLink("A", "C")
	g.AddLink("B", "A")
	g.AddLink("C", "A")
	g.AddLink("D", "C") // 测试悬挂节点
	pr := CalculatePageRank(g, 0.85, 20)

	// 验证PR值总和≈1（允许浮点误差）
	total := 0.0
	for _, v := range pr {
		total += v
	}
	if math.Abs(total-1.0) > 1e-10 {
		t.Errorf("PR sum should be 1.0, got %f", total)
	}
	// 验证A的PR值应最高（被B/C同时引用）
	if pr["A"] <= pr["B"] || pr["A"] <= pr["C"] {
		t.Errorf("A should have highest PR: %v", pr)
	}
	// 验证悬挂节点D的PR值分配
	if pr["D"] <= 0 {
		t.Errorf("Dangling node D should have PR > 0: %f", pr["D"])
	}
}

func TestSingleNode(t *testing.T) {
	g := NewGraph()
	g.AddLink("A", "A") // 自环

	pr := CalculatePageRank(g, 0.85, 10)
	if math.Abs(pr["A"]-1.0) > 1e-10 {
		t.Errorf("Single node PR should be 1.0, got %f", pr["A"])
	}
}
