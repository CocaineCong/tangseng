package trie

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (trie *Trie) FindAllByPrefixForRecall(prefix string) []string {
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := string(prefixs[i])
		if _, ok := node.ChildrenRecall[c]; !ok {
			return nil
		}
		node, _ = node.ChildrenRecall[c] // nolint:golint,gosimple
	}
	words := make([]string, 0)
	trie.dfsForRecall(node, prefix, &words)
	return words
}

func (trie *Trie) dfsForRecall(node *TrieNode, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}

	for c, child := range node.ChildrenRecall {
		trie.dfsForRecall(child, word+c, words)
	}
}

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

// ParseTrieNode 解析 TrieNode 结构体
func ParseTrieNode(str string) (*TrieNode, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}

	node := &TrieNode{
		IsEnd:          false,
		ChildrenRecall: make(map[string]*TrieNode),
	}

	for key, value := range data {
		childData, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid child data for key: %s", key)
		}

		childNode, err := ParseTrieNodeJSON(childData)
		if err != nil {
			return nil, err
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}

// 解析 TrieNode 结构体的 JSON 数据
func ParseTrieNodeJSON(data map[string]interface{}) (*TrieNode, error) {
	node := &TrieNode{
		IsEnd:          false,
		ChildrenRecall: make(map[string]*TrieNode),
	}

	isEnd, ok := data["is_end"].(bool)
	if ok {
		node.IsEnd = isEnd
	}

	childrenData, ok := data["children_recall"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid children data")
	}

	for key, value := range childrenData {
		childData, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid child data for key: %s", key)
		}

		childNode, err := ParseTrieNodeJSON(childData)
		if err != nil {
			return nil, err
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}
