package storage

import (
	"bytes"
	"context"
	"os"

	"github.com/pkg/errors"

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
		err = errors.Wrap(err, "failed marshal data")
		return
	}
	if err = d.PutTrieTree([]byte(consts.TrieTreeBucket), tt); err != nil {
		err = errors.Wrap(err, "failed to put trie tree")
	}
	return
}

// GetTrieTreeDict 获取 trie tree
func (d *TrieDB) GetTrieTreeDict() (trieTree *trie.Trie, err error) {
	v, err := d.GetTrieTree([]byte(consts.TrieTreeBucket))
	if err != nil {
		err = errors.WithMessage(err, "getTrieTree error")
		return
	}
	replaced := bytes.Replace(v, []byte("children"), []byte("children_recall"), -1)
	node, err := trie.ParseTrieNode(string(replaced))
	if err != nil {
		err = errors.WithMessage(err, "ParseTrieNode error")
		return
	}

	trieTree = trie.NewTrie()
	trieTree.Root = node

	return
}

// PutTrieTree 存储
func (d *TrieDB) PutTrieTree(key, value []byte) (err error) {
	err = Put(d.db, consts.TrieTreeBucket, key, value)
	if err != nil {
		err = errors.WithMessage(err, "put error")
	}
	return
}

// GetTrieTree 通过term获取value
func (d *TrieDB) GetTrieTree(key []byte) (value []byte, err error) {
	value, err = Get(d.db, consts.TrieTreeBucket, key)
	err = errors.WithMessage(err, "get error")
	return
}

// Close 关闭db
func (d *TrieDB) Close() (err error) {
	if err = d.db.Close(); err != nil {
		err = errors.WithMessage(d.db.Close(), "close error")
	}
	return
}
