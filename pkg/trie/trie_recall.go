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
	"encoding/json"

	"github.com/pkg/errors"
)

// FindAllByPrefixForRecall 召回专用的，通过前缀获取所有的节点
func (trie *Trie) FindAllByPrefixForRecall(prefix string) []string {
	prefixes := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixes); i++ {
		c := string(prefixes[i])
		if _, ok := node.ChildrenRecall[c]; !ok {
			return nil
		}
		node, _ = node.ChildrenRecall[c] // nolint:golint,gosimple
	}
	words := make([]string, 0)
	trie.dfsForRecall(node, prefix, &words)
	return words
}

func (trie *Trie) dfsForRecall(node *Node, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}

	for c, child := range node.ChildrenRecall {
		trie.dfsForRecall(child, word+c, words)
	}
}

// SearchForRecall 召回时查询这个word是否存在
func (trie *Trie) SearchForRecall(word string) bool {
	words := []rune(word)
	node := trie.Root
	for i := 0; i < len(words); i++ {
		c := string(words[i])
		if _, ok := node.ChildrenRecall[c]; !ok {
			return false
		}
		node, _ = node.ChildrenRecall[c] // nolint:golint,gosimple
	}
	return node.IsEnd
}

// ParseTrieNode 解析 TrieNode 结构体,
// 从数据库读出来是字符串，然后解析这个字符串成一棵树
func ParseTrieNode(str string) (*Node, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal data")
	}

	node := &Node{
		IsEnd:          false,
		ChildrenRecall: make(map[string]*Node),
	}

	for key, value := range data {
		childData, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.Wrap(errors.Errorf("invalid child data for key: %s", key), "failed to assert type")
		}

		childNode, err := parseTrieNodeChild(childData)
		if err != nil {
			return nil, errors.WithMessage(err, "parseTrieNodeChild error")
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}

// parseTrieNodeChild 解析 TrieNode 结构体的 JSON 数据
// 将这个map数据解析成树状数据
func parseTrieNodeChild(data map[string]interface{}) (*Node, error) {
	node := &Node{
		IsEnd:          false,
		ChildrenRecall: make(map[string]*Node),
	}

	isEnd, ok := data["is_end"].(bool)
	if ok {
		node.IsEnd = isEnd
	}

	childrenData, ok := data["children_recall"].(map[string]interface{})
	if !ok {
		return nil, errors.Wrap(errors.New("invalid children data"), "failed to assert type")
	}

	for key, value := range childrenData {
		childData, ok := value.(map[string]interface{})
		if !ok {
			return nil, errors.Wrap(errors.Errorf("invalid child data for key: %s", key), "failed to assert type")
		}

		childNode, err := parseTrieNodeChild(childData)
		if err != nil {
			return nil, errors.WithMessage(err, "parseTrieNodeChild error")
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}

// TraverseForRecall 查看所有节点的信息
func (trie *Trie) TraverseForRecall() {
	traverseForRecall(trie.Root, "")
}

func traverseForRecall(node *Node, prefix string) {
	if node.IsEnd {
		return
	}

	for c, child := range node.ChildrenRecall {
		traverseForRecall(child, prefix+c)
	}
}
