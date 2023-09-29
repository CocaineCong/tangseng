package storage

import (
	"fmt"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
)

type TrieDB struct {
	db *bolt.DB
}

var GlobalTrieDBs *TrieDB

const InvertedDBPaths = "/Users/mac/GolandProjects/Go-SearchEngine/app/index_platform/trie_data/"

// InitTrieDBs 初始化trie
func InitTrieDBs() { // TODO: 先都放在一个下面吧，后面再lb到多个文件
	filePath := fmt.Sprintf("%sinit.%s", InvertedDBPaths, consts.TrieTreeBucket)
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	GlobalTrieDBs = &TrieDB{db}
}

// NewTrieDB 新建一个forward db对象
func NewTrieDB(dbName string) (*TrieDB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.LogrusObj.Errorf("NewTrieDB: %+v", err)
		return nil, err
	}

	return &TrieDB{db}, nil
}

func (d *TrieDB) StorageDict(trieTree *trie.Trie) (err error) {
	trieByte, _ := trieTree.MarshalJSON()
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
