package storage

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"testing"

	"github.com/CocaineCong/tangseng/pkg/trie"
)

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
	if err != nil {
		fmt.Println(err)
	}
	tree2.Traverse()
}

func TestBinaryCnTree(t *testing.T) {
	// 假设这是要编码的中文字符串
	data := []byte("你好，世界！")
	// 假设这是要压缩的二进制数据

	// 压缩数据
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, _ = gzWriter.Write(data)
	_ = gzWriter.Close()
	fmt.Println(buf.Bytes())
	// 解压数据
	gzReader, err := gzip.NewReader(&buf)
	if err != nil {
		panic(err)
	}
	defer func(gzReader *gzip.Reader) {
		err = gzReader.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(gzReader)

	result, err := io.ReadAll(gzReader)
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
	_, _ = gzWriter.Write(data)
	_ = gzWriter.Close()
	fmt.Println(buf.Bytes())
	// 解压数据
	gzReader, err := gzip.NewReader(&buf)
	if err != nil {
		panic(err)
	}
	defer func(gzReader *gzip.Reader) {
		err = gzReader.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(gzReader)

	result, err := io.ReadAll(gzReader)
	if err != nil {
		panic(err)
	}

	tt2 := trie.NewTrie()
	err = tt2.UnmarshalJSON(result)
	if err != nil {
		fmt.Println(err)
	}
	// 打印解压后的二进制数据
	tt2.Traverse()
}
