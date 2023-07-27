package storage

import (
	"encoding/json"

	"github.com/spf13/cast"
	bolt "go.etcd.io/bbolt"

	log "github.com/CocaineCong/tangseng/pkg/logger"
)

const forwardBucket = "forward"

const ForwardCountKey = "forwardCount"

type ForwardDB struct {
	db *bolt.DB
}

// NewForwardDB 新建一个forward db对象
func NewForwardDB(dbName string) (*ForwardDB, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.LogrusObj.Errorf("NewForwardDB: %v", err.Error())
		return nil, err
	}

	return &ForwardDB{db}, nil
}

// AddForwardByDoc 通过doc进行存储
func (f *ForwardDB) AddForwardByDoc(doc *Document) error {
	key := cast.ToString(doc.DocId)
	body, _ := json.Marshal(doc)
	return Put(f.db, forwardBucket, []byte(key), body)
}

// PutForwardByKV 通过kv进行存储
func (f *ForwardDB) PutForwardByKV(key, value []byte) error {
	return Put(f.db, forwardBucket, key, value)
}

// ForwardCount 获取文档总数
func (f *ForwardDB) ForwardCount() (r int64, err error) {
	body, err := Get(f.db, forwardBucket, []byte(ForwardCountKey))
	if err != nil {
		return
	}

	r = cast.ToInt64(body)
	return
}

// UpdateForwardCount 获取文档总数
func (f *ForwardDB) UpdateForwardCount(count int64) error {
	return Put(f.db, forwardBucket, []byte(ForwardCountKey), []byte(cast.ToString(count)))
}

// GetForward 获取forward数据
func (f *ForwardDB) GetForward(docId int64) (r []byte, err error) {
	return Get(f.db, forwardBucket, []byte(cast.ToString(docId)))
}

// GetForwardCursor 获取遍历游标
func (f *ForwardDB) GetForwardCursor(termCh chan KvInfo) error {
	return f.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(forwardBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			termCh <- KvInfo{k, v}
		}
		close(termCh)
		return nil
	})
}

// Close 关闭db
func (f *ForwardDB) Close() error {
	return f.db.Close()
}
