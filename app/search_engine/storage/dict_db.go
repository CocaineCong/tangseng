package storage

import (
	"bytes"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/trie"
	"github.com/CocaineCong/tangseng/pkg/util/codec"
)

type DictDB struct {
	db *bolt.DB
}

// NewDictDB 新建一个forward db对象
func NewDictDB(dbName string) (*DictDB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.LogrusObj.Errorf("NewDictDB: %+v", err)
		return nil, err
	}

	return &DictDB{db}, nil
}

func (d *DictDB) StorageDict(trieTree *trie.Trie) (err error) {
	buf := bytes.NewBuffer(nil)
	err = codec.BinaryEncoding(buf, trieTree)
	if err != nil {
		return
	}

	err = d.PutTrieTree([]byte(consts.DictBucket), buf.Bytes())

	return
}

// GetTrieTreeDict 获取 trie tree
func (d *DictDB) GetTrieTreeDict(buf *bytes.Buffer, trieTree *trie.Trie) (err error) {
	v, err := d.GetTrieTree([]byte(consts.DictBucket))
	if err != nil {
		return
	}
	buf = bytes.NewBuffer(v)
	err = codec.BinaryDecoding(buf, trieTree)
	if err != nil {
		return
	}

	return
}

// PutTrieTree 存储
func (d *DictDB) PutTrieTree(key, value []byte) error {
	return Put(d.db, consts.DictBucket, key, value)
}

// GetTrieTree 通过term获取value
func (d *DictDB) GetTrieTree(key []byte) (value []byte, err error) {
	return Get(d.db, consts.DictBucket, key)
}

// Close 关闭db
func (d *DictDB) Close() error {
	return d.db.Close()
}
