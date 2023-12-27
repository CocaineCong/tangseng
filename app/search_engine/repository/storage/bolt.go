package storage

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return errors.Wrap(err, "failed to create bucket")
		}
		err = b.Put(key, value)
		if err != nil {
			err = errors.Wrap(err, "failed to put data")
		}
		return err
	})
}

// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) (r []byte, err error) {
	err = db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			b, _ = tx.CreateBucketIfNotExists([]byte(bucket))
		}
		r = b.Get(key)
		if r == nil { // 如果是空的话，直接创建这个key，然后返回这个key的初始值，也就是0
			r = []byte("0")
			return
		}
		return
	})
	if err != nil {
		err = errors.Wrap(err, "view error")
	}
	return
}
