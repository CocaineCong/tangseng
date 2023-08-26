package storage

import (
	"bytes"

	"github.com/spf13/cast"
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

func (d *DictDB) StorageDict(segId int, trieTree *trie.Trie) (err error) {
	buf := bytes.NewBuffer(nil)
	err = codec.BinaryEncoding(buf, trieTree)
	if err != nil {
		return
	}

	err = d.PutTrimTreeByKV([]byte(cast.ToString(segId)), buf.Bytes())

	return
}

// PutTrimTreeByKV 通过kv进行存储
func (d *DictDB) PutTrimTreeByKV(key, value []byte) error {
	return Put(d.db, consts.DictBucket, key, value)
}

// GetTrimTree 通过term获取value
func (d *DictDB) GetTrimTree(key []byte) (value []byte, err error) {
	return Get(d.db, consts.DictBucket, key)
}

// Close 关闭db
func (d *DictDB) Close() error {
	return d.db.Close()
}
