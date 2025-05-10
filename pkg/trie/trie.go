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

package trie

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// Node TODO:后面看看能不能把build和recall的过程分开,主要是 cmap.ConcurrentMap[string, *TrieNode] 没法反序列化...
type Node struct {
	IsEnd          bool                              `json:"is_end"`   // 标记该节点是否为一个单词的末尾
	Children       cmap.ConcurrentMap[string, *Node] `json:"children"` // 存储子节点的指针
	ChildrenRecall map[string]*Node                  `json:"children_recall"`
}

func NewTrieNode() *Node {
	m := cmap.New[*Node]()
	return &Node{
		IsEnd:    false,
		Children: m,
	}
}

type Trie struct {
	Root *Node `json:"root"` // 存储 Trie 树的根节点
}

func NewTrie() *Trie {
	return &Trie{Root: NewTrieNode()}
}

func (trie *Trie) Insert(word string) {
	words := []rune(word)
	node := trie.Root
	for i := 0; i < len(words); i++ {
		c := string(words[i])
		if _, ok := node.Children.Get(c); !ok {
			node.Children.Set(c, NewTrieNode())
		}
		node, _ = node.Children.Get(c)
	}
	node.IsEnd = true
}

func (trie *Trie) Search(word string) bool {
	words := []rune(word)
	node := trie.Root
	for i := 0; i < len(words); i++ {
		c := string(words[i])
		if _, ok := node.Children.Get(c); !ok {
			return false
		}
		node, _ = node.Children.Get(c)
	}
	return node.IsEnd
}

func (trie *Trie) StartsWith(prefix string) bool {
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := string(prefixs[i])
		if _, ok := node.Children.Get(c); !ok {
			return false
		}
		node, _ = node.Children.Get(c)
	}
	return true
}

func (trie *Trie) FindAllByPrefix(prefix string) []string {
	prefixes := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixes); i++ {
		c := string(prefixes[i])
		if _, ok := node.Children.Get(c); !ok {
			return nil
		}
		node, _ = node.Children.Get(c)
	}
	words := make([]string, 0)
	trie.dfs(node, prefix, &words)
	return words
}

func (trie *Trie) dfs(node *Node, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}

	for c, child := range node.Children.Items() {
		trie.dfs(child, word+c, words)
	}
}

func (trie *Trie) Merge(other *Trie) {
	if other == nil {
		return
	}

	var mergeNodes func(n1, n2 *Node)
	mergeNodes = func(n1, n2 *Node) {
		for c, child := range n2.Children.Items() {
			if v, ok := n1.Children.Get(c); ok {
				mergeNodes(v, child)
			} else {
				n1.Children.Set(c, child)
			}
		}
		n1.IsEnd = n1.IsEnd || n2.IsEnd
	}

	mergeNodes(trie.Root, other.Root)
}

func traverse(node *Node, prefix string) {
	if node.IsEnd {
		fmt.Println(prefix)
	}

	for c, child := range node.Children.Items() {
		traverse(child, prefix+c)
	}
}

func (trie *Trie) Traverse() {
	traverse(trie.Root, "")
}
