package trie

import (
	"encoding/json"
	"errors"
	"fmt"
)

// FindAllByPrefixForRecall 召回专用的，通过前缀获取所有的节点
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

		childNode, err := parseTrieNodeChild(childData)
		if err != nil {
			return nil, err
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}

// parseTrieNodeChild 解析 TrieNode 结构体的 JSON 数据
// 将这个map数据解析成树状数据
func parseTrieNodeChild(data map[string]interface{}) (*TrieNode, error) {
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

		childNode, err := parseTrieNodeChild(childData)
		if err != nil {
			return nil, err
		}

		node.ChildrenRecall[key] = childNode
	}

	return node, nil
}

// TraverseForRecall 查看所有节点的信息
func (trie *Trie) TraverseForRecall() {
	traverseForRecall(trie.Root, "")
}

func traverseForRecall(node *TrieNode, prefix string) {
	if node.IsEnd {
		fmt.Println(prefix)
	}

	for c, child := range node.ChildrenRecall {
		traverseForRecall(child, prefix+c)
	}
}
