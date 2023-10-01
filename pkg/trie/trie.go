package trie

import (
	"encoding/json"
	"errors"
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// TrieNode TODO:后面看看能不能把build和recall的过程分开,主要是 cmap.ConcurrentMap[string, *TrieNode] 没法反序列化...
type TrieNode struct {
	IsEnd          bool                                  `json:"is_end"`   // 标记该节点是否为一个单词的末尾
	Children       cmap.ConcurrentMap[string, *TrieNode] `json:"children"` // 存储子节点的指针
	ChildrenRecall map[string]*TrieNode                  `json:"children_recall"`
}

func NewTrieNode() *TrieNode {
	m := cmap.New[*TrieNode]()
	return &TrieNode{
		IsEnd:    false,
		Children: m,
	}
}

type Trie struct {
	Root *TrieNode `json:"root"` // 存储 Trie 树的根节点
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
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := string(prefixs[i])
		if _, ok := node.Children.Get(c); !ok {
			return nil
		}
		node, _ = node.Children.Get(c)
	}
	words := make([]string, 0)
	trie.dfs(node, prefix, &words)
	return words
}

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

func (trie *Trie) dfs(node *TrieNode, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}

	for c, child := range node.Children.Items() {
		trie.dfs(child, word+c, words)
	}
}

func (trie *Trie) dfsForRecall(node *TrieNode, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}

	for c, child := range node.ChildrenRecall {
		trie.dfsForRecall(child, word+c, words)
	}
}

func (trie *Trie) Merge(other *Trie) {
	if other == nil {
		return
	}

	var mergeNodes func(n1, n2 *TrieNode)
	mergeNodes = func(n1, n2 *TrieNode) {
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

func traverse(node *TrieNode, prefix string) {
	if node.IsEnd {
		fmt.Println(prefix)
	}

	for c, child := range node.Children.Items() {
		traverse(child, prefix+c)
	}
}

func traverseForRecall(node *TrieNode, prefix string) {
	if node.IsEnd {
		fmt.Println(prefix)
	}

	for c, child := range node.ChildrenRecall {
		traverseForRecall(child, prefix+c)
	}
}

func (trie *Trie) TraverseForRecall() {
	traverseForRecall(trie.Root, "")
}

func (trie *Trie) Traverse() {
	traverse(trie.Root, "")
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
