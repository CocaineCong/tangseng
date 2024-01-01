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
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

type TrieDB struct {
	file *os.File
	db   *bolt.DB
}

// NewTrieDB 初始化trie
func NewTrieDB(filePath string) *TrieDB { // TODO: 先都放在一个下面吧，后面再lb到多个文件
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.LogrusObj.Error(err)
	}

	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil
	}

	return &TrieDB{f, db}
}

func (d *TrieDB) StorageDict(trieTree *trie.Trie) (err error) {
	trieByte, _ := trieTree.Root.Children.MarshalJSON()
	err = d.PutTrieTree([]byte(consts.TrieTreeBucket), trieByte)

	return
}

// GetTrieTreeInfo 获取 trie tree
func (d *TrieDB) GetTrieTreeInfo() (trieTree *trie.Trie, err error) {
	v, err := d.GetTrieTree([]byte(consts.TrieTreeBucket))
	if err != nil {
		return
	}

	trieTree = trie.NewTrie()
	err = trieTree.UnmarshalJSON(v)

	return
}

// PutTrieTree 存储
func (d *TrieDB) PutTrieTree(key, value []byte) error {
	return Put(d.db, consts.TrieTreeBucket, key, value)
}

// GetTrieTree 通过term获取value
func (d *TrieDB) GetTrieTree(key []byte) (value []byte, err error) {
	return Get(d.db, consts.TrieTreeBucket, key)
}

// Close 关闭db
func (d *TrieDB) Close() error {
	return d.db.Close()
}
