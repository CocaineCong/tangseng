package storage

import (
	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) (r []byte, err error) {
	err = db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(bucket))
		r = b.Get(key)
		if r == nil { // 如果是空的话，直接创建这个key，然后返回这个key的初始值，也就是0
			_, err = b.CreateBucketIfNotExists(key)
			if err != nil {
				return
			}
			err = b.Put(key, []byte(consts.ForwardCountInitValue))
			if err != nil {
				return
			}
			r = b.Get(key)
			return
		}
		return
	})

	return
}
