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

type Graph struct {
	Nodes   []string
	OutLink map[string][]string // 节点指向的其他节点
	InLink  map[string][]string // 指向节点的其他节点
}

func NewGraph() *Graph {
	return &Graph{
		OutLink: make(map[string][]string),
		InLink:  make(map[string][]string),
	}
}

func (g *Graph) AddLink(from, to string) {
	g.OutLink[from] = append(g.OutLink[from], to)
	g.InLink[to] = append(g.InLink[to], from)
	if !g.containsNode(from) {
		g.Nodes = append(g.Nodes, from)
	}
	if !g.containsNode(to) {
		g.Nodes = append(g.Nodes, to)
	}
}

func (g *Graph) containsNode(node string) bool {
	for _, n := range g.Nodes {
		if n == node {
			return true
		}
	}
	return false
}

func CalculatePageRank(g *Graph, damping float64, iterations int) map[string]float64 {
	N := float64(len(g.Nodes))
	pr := make(map[string]float64)

	// 初始化PR值
	for _, node := range g.Nodes {
		pr[node] = 1.0 / N
	}

	for i := 0; i < iterations; i++ {
		newPR := make(map[string]float64)
		leak := 0.0
		// 计算每个节点的新PR值
		for _, node := range g.Nodes {
			sum := 0.0
			for _, inNode := range g.InLink[node] {
				sum += pr[inNode] / float64(len(g.OutLink[inNode]))
			}
			newPR[node] = (1-damping)/N + damping*sum
		}
		// 处理悬挂节点（无出链）
		for _, node := range g.Nodes {
			if len(g.OutLink[node]) == 0 {
				leak += damping * pr[node] / N
			}
		}
		// 分配泄漏的PR值
		for node := range newPR {
			newPR[node] += leak
		}
		pr = newPR
	}

	return pr
}
