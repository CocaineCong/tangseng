package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/pkg/fileutils"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

func TestTrieDB_GetTrieTree(t *testing.T) {
	aConfig := config.Conf.SeConfig.StoragePath + "0.Trie"
	d, _ := NewTrieDB(aConfig)
	trieTree := trie.NewTrie()
	trieTree, err := d.GetTrieTreeInfo()
	fmt.Println(err)
	// trieTree.Traverse()
	a := trieTree.FindAllByPrefix("传")
	fmt.Println(a)
}

func TestBinaryTrieTree(t *testing.T) {
	tree := trie.NewTrie()
	tree.Insert("你好")
	tree.Traverse()
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(tree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("tree2")
	tree2 := trie.NewTrie()
	err = gob.NewDecoder(buf).Decode(tree2)
	tree2.Traverse()
}

func TestBinaryCnTree(t *testing.T) {
	// 假设这是要编码的中文字符串
	data := []byte("你好，世界！")
	// 假设这是要压缩的二进制数据

	// 压缩数据
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	gzWriter.Write(data)
	gzWriter.Close()
	fmt.Println(buf.Bytes())
	// 解压数据
	gzReader, err := gzip.NewReader(&buf)
	if err != nil {
		panic(err)
	}
	defer gzReader.Close()

	result, err := ioutil.ReadAll(gzReader)
	if err != nil {
		panic(err)
	}

	// 打印解压后的二进制数据
	fmt.Printf("%s", result)

}

func TestBinaryCnStrTree(t *testing.T) {
	// 假设这是要编码的中文字符串
	tt := trie.NewTrie()
	tt.Insert("你好")
	tt.Insert("世界")

	data, _ := tt.MarshalJSON()
	// 假设这是要压缩的二进制数据

	// 压缩数据
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	gzWriter.Write(data)
	gzWriter.Close()
	fmt.Println(buf.Bytes())
	// 解压数据
	gzReader, err := gzip.NewReader(&buf)
	if err != nil {
		panic(err)
	}
	defer gzReader.Close()

	result, err := io.ReadAll(gzReader)
	if err != nil {
		panic(err)
	}

	tt2 := trie.NewTrie()
	err = tt2.UnmarshalJSON(result)

	// 打印解压后的二进制数据
	tt2.Traverse()
}

func TestGetFiles(t *testing.T) {
	filePath := fileutils.GetFiles(InvertedDBPaths)
	fmt.Println(filePath)
}

func TestInitTrieDBs(t *testing.T) {
	InitTrieDBs()
	fmt.Println(GlobalTrieDBs)
}

func TestTrieDB_GetTrieTreeInfo(t *testing.T) {
	InitTrieDBs()
	tree, err := GlobalTrieDBs[0].GetTrieTreeInfo()
	if err != nil {
		fmt.Println(err)
	}
	tree.Traverse()
}
