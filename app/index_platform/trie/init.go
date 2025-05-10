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

package trie

import (
	"github.com/CocaineCong/tangseng/pkg/trie"
)

var GlobalTrieTree *trie.Trie

func InitTrieTree() {
	// GlobalTrieTree = trie.NewTrie()
	// TODO: 这里的想法是把原始的读出来合并的，但是第一次读的时候由于是空的，所以会强制报错，用recover也不起作用，后面看看怎么处理吧... :-(
	// val, err := storage.GlobalTrieDBs.GetTrieTreeInfo()
	// if err != nil {
	// 	// 第一次读取会出现没有的情况
	// } else {
	// 	GlobalTrieTree.Merge(val)
	// }
}
