package storage

import (
	"bytes"
	"context"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
	"github.com/CocaineCong/tangseng/repository/redis"
)

var GlobalTrieDB []*TrieDB

type TrieDB struct {
	file *os.File
	db   *bolt.DB
}

// InitGlobalTrieDB 初始化trie tree树
func InitGlobalTrieDB(ctx context.Context) {
	dbs := make([]*TrieDB, 0)
	filePath, _ := redis.ListInvertedPath(ctx, redis.TireTreeDbPathKey)
	for _, file := range filePath {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.LogrusObj.Error(err)
		}

		db, err := bolt.Open(file, 0600, nil)
		if err != nil {
			log.LogrusObj.Error(err)
		}
		dbs = append(dbs, &TrieDB{f, db})
	}
	if len(filePath) == 0 {
		return
	}
	GlobalTrieDB = dbs
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
	tt, err := trieTree.MarshalJSON()
	if err != nil {
		return
	}

	err = d.PutTrieTree([]byte(consts.TrieTreeBucket), tt)

	return
}

// GetTrieTreeDict 获取 trie tree
func (d *TrieDB) GetTrieTreeDict() (trieTree *trie.Trie, err error) {
	v, err := d.GetTrieTree([]byte(consts.TrieTreeBucket))
	if err != nil {
		return
	}
	replaced := bytes.Replace(v, []byte("children"), []byte("children_recall"), -1)
	node, err := trie.ParseTrieNode(string(replaced))
	if err != nil {
		return
	}

	trieTree = trie.NewTrie()
	trieTree.Root = node

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
