package trie

import (
	"fmt"
	"sync"
)

type TrieNode struct {
	IsEnd bool `json:"is_end"` // 标记该节点是否为一个单词的末尾
	// Children map[rune]*TrieNode `json:"children"` // 存储子节点的指针
	Children *sync.Map `json:"children"` // 存储子节点的指针

}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		IsEnd: false,
		// Children: make(map[rune]*TrieNode),
		Children: new(sync.Map),
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
		c := words[i]
		// if _, ok := node.Children[c]; !ok {
		if _, ok := node.Children.Load(c); !ok {
			node.Children.Store(c, NewTrieNode())
		}
		nodeTmp, _ := node.Children.Load(c)
		node = nodeTmp.(*TrieNode)
	}
	node.IsEnd = true
}

func (trie *Trie) Search(word string) bool {
	words := []rune(word)
	node := trie.Root
	for i := 0; i < len(words); i++ {
		c := words[i]
		if _, ok := node.Children.Load(c); !ok {
			return false
		}
		nodeTmp, _ := node.Children.Load(c)
		node = nodeTmp.(*TrieNode)
	}
	return node.IsEnd
}

func (trie *Trie) StartsWith(prefix string) bool {
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := prefixs[i]
		if _, ok := node.Children.Load(c); !ok {
			return false
		}
		nodeTmp, _ := node.Children.Load(c)
		node = nodeTmp.(*TrieNode)
	}
	return true
}

func (trie *Trie) FindAllByPrefix(prefix string) []string {
	prefixs := []rune(prefix)
	node := trie.Root
	for i := 0; i < len(prefixs); i++ {
		c := prefixs[i]
		if _, ok := node.Children.Load(c); !ok {
			return nil
		}
		nodeTmp, _ := node.Children.Load(c)
		node = nodeTmp.(*TrieNode)
	}
	words := make([]string, 0)
	trie.dfs(node, prefix, &words)
	return words
}

func (trie *Trie) dfs(node *TrieNode, word string, words *[]string) {
	if node.IsEnd {
		*words = append(*words, word)
	}
	node.Children.Range(func(key, value any) bool {
		trie.dfs(value.(*TrieNode), word+string(key.(rune)), words)
		return true
	})

	// for c, child := range node.Children {
	// 	trie.dfs(child, word+string(c), words)
	// }
}

func (trie *Trie) Merge(other *Trie) {
	if other == nil {
		return
	}

	var mergeNodes func(n1, n2 *TrieNode)
	mergeNodes = func(n1, n2 *TrieNode) {
		// for c, child := range n2.Children {
		// 	if _, ok := n1.Children[c]; ok {
		// 		mergeNodes(n1.Children[c], child)
		// 	} else {
		// 		n1.Children[c] = child
		// 	}
		// }
		n2.Children.Range(func(key, value any) bool {
			if val, ok := n1.Children.Load(key); ok {
				mergeNodes(val.(*TrieNode), value.(*TrieNode))
			} else {
				n1.Children.Store(key, value)
			}

			return true
		})

		n1.IsEnd = n1.IsEnd || n2.IsEnd
	}

	mergeNodes(trie.Root, other.Root)
}

func traverse(node *TrieNode, prefix string) {
	if node.IsEnd {
		fmt.Println(prefix)
	}

	// for c, child := range node.Children {
	// 	traverse(child, prefix+string(c))
	// }
	node.Children.Range(func(key, value any) bool {
		traverse(value.(*TrieNode), prefix+string(key.(rune)))

		return true
	})
}

func (trie *Trie) Traverse() {
	traverse(trie.Root, "")
}
