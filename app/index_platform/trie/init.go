package trie

import (
	"github.com/CocaineCong/tangseng/pkg/trie"
)

var GobalTrieTree *trie.Trie

func InitTrieTree() {
	// GobalTrieTree = trie.NewTrie()
	// TODO: 这里的想法是把原始的读出来合并的，但是第一次读的时候由于是空的，所以会强制报错，用recover也不起作用，后面看看怎么处理吧... :-(
	// val, err := storage.GlobalTrieDBs.GetTrieTreeInfo()
	// if err != nil {
	// 	// 第一次读取会出现没有的情况
	// } else {
	// 	GobalTrieTree.Merge(val)
	// }
}
