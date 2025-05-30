// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
