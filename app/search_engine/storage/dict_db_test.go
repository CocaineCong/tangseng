package storage

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

func TestDictDB_GetTrimTree(t *testing.T) {
	aConfig := config.Conf.SeConfig.StoragePath + "0.dict"
	d, _ := NewDictDB(aConfig)
	buf := bytes.NewBuffer(nil)
	trieTree := trie.NewTrie()
	err := d.GetTrieTreeDict(buf, trieTree)
	fmt.Println(err)
	a := trieTree.Find("导")
	fmt.Println(a)
}
