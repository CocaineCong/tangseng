# 前缀树相关

```go
type TrieNode struct {
	IsEnd          bool                                  `json:"is_end"`   // 标记该节点是否为一个单词的末尾
	Children       cmap.ConcurrentMap[string, *TrieNode] `json:"children"` // 存储子节点的指针
	ChildrenRecall map[string]*TrieNode                  `json:"children_recall"`
}
```

插入前缀树

```go
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
```

查询前缀树

```go
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
```