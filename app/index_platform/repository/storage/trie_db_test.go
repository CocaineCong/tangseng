// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
