package storage

import (
	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

type DictDB struct {
	db *bolt.DB
}

// NewDictDB 新建一个forward db对象
func NewDictDB(dbName string) (*DictDB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.LogrusObj.Errorf("NewDictDB: %v", err.Error())
		return nil, err
	}

	return &DictDB{db}, nil
}

// PutForwardByKV 通过kv进行存储
func (d *DictDB) PutForwardByKV(key, value []byte) error {
	return Put(d.db, consts.DictBucket, key, value)
}
