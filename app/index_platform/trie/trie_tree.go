package trie

import (
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

var GobalTrieTree *trie.Trie

func InitTrieTree() {
	GobalTrieTree = trie.NewTrie()
	for _, trieTree := range storage.GobalTrieDBs {
		val, _ := trieTree.GetTrieTreeInfo()
		GobalTrieTree.Merge(val)
	}
}
