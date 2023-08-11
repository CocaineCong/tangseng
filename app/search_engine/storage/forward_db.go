package storage

import (
	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

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
	body, _ := sonic.Marshal(doc.Body)
	return Put(f.db, consts.ForwardBucket, []byte(key), body)
}

// PutForwardByKV 通过kv进行存储
func (f *ForwardDB) PutForwardByKV(key, value []byte) error {
	return Put(f.db, consts.ForwardBucket, key, value)
}

// ForwardCount 获取文档总数
func (f *ForwardDB) ForwardCount() (r int64, err error) {
	body, err := Get(f.db, consts.ForwardBucket, []byte(consts.ForwardCountKey))
	if err != nil {
		return
	}

	r = cast.ToInt64(body)
	return
}

// UpdateForwardCount 获取文档总数
func (f *ForwardDB) UpdateForwardCount(count int64) error {
	return Put(f.db, consts.ForwardBucket, []byte(consts.ForwardCountKey), []byte(cast.ToString(count)))
}

// GetForward 获取forward数据
func (f *ForwardDB) GetForward(docId int64) (r []byte, err error) {
	return Get(f.db, consts.ForwardBucket, []byte(cast.ToString(docId)))
}

// GetForwardCursor 获取遍历游标
func (f *ForwardDB) GetForwardCursor(termCh chan KvInfo) error {
	return f.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(consts.ForwardBucket))
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
