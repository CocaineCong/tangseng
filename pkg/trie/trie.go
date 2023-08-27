package trie

import (
	"fmt"
)

type TrieNode struct {
	IsEnd    bool               `json:"is_end"`   // 标记该节点是否为一个单词的末尾
	Children map[byte]*TrieNode `json:"children"` // 存储子节点的指针
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		IsEnd:    false,
		Children: make(map[byte]*TrieNode),
	}
}

type Trie struct {
	Root *TrieNode // 存储 Trie 树的根节点
}

func NewTrie() *Trie {
	return &Trie{Root: NewTrieNode()}
}

func (trie *Trie) Insert(word string) {
	node := trie.Root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if _, ok := node.Children[c]; !ok {
			node.Children[c] = NewTrieNode()
		}
		node = node.Children[c]
	}
	node.IsEnd = true
}

func (trie *Trie) Search(word string) bool {
	node := trie.Root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if _, ok := node.Children[c]; !ok {
			return false
		}
		node = node.Children[c]
	}
	return node.IsEnd
}

func (trie *Trie) StartsWith(prefix string) bool {
	node := trie.Root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		if _, ok := node.Children[c]; !ok {
			return false
		}
		node = node.Children[c]
	}
	return true
}

func (trie *Trie) FindAllByPrefix(prefix string) []string {
	node := trie.Root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		if _, ok := node.Children[c]; !ok {
			return nil
		}
		node = node.Children[c]
	}
	words := make([]string, 0)
	trie.dfs(node, prefix, &words)
	return words
}

func (trie *Trie) dfs(node *TrieNode, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}
	for c, child := range node.Children {
		trie.dfs(child, word+string(c), words)
	}
}

func (trie *Trie) Merge(other *Trie) {
	if other == nil {
		return
	}

	var mergeNodes func(n1, n2 *TrieNode)
	mergeNodes = func(n1, n2 *TrieNode) {
		for c, child := range n2.Children {
			if _, ok := n1.Children[c]; ok {
				mergeNodes(n1.Children[c], child)
			} else {
				n1.Children[c] = child
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

	for c, child := range node.Children {
		traverse(child, prefix+string(c))
	}
}

func (trie *Trie) Traverse() {
	traverse(trie.Root, "")
}
