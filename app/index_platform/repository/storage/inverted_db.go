package storage

import (
	"github.com/pkg/errors"
	"os"

	"github.com/RoaringBitmap/roaring"
	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

type KvInfo struct {
	Key   []byte
	Value []byte
}

type InvertedDB struct { // 后续做mmap(这个好难)
	file   *os.File
	db     *bolt.DB
	offset int64
}

// NewInvertedDB 新建一个inverted
func NewInvertedDB(invertedName string) *InvertedDB {
	f, err := os.OpenFile(invertedName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.LogrusObj.Error(err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.LogrusObj.Error(err)
	}
	db, err := bolt.Open(invertedName, 0600, nil)
	if err != nil {
		log.LogrusObj.Error(err)
	}

	return &InvertedDB{f, db, stat.Size()}
}

// StoragePostings 存储 倒排索引表
func (t *InvertedDB) StoragePostings(token string, values []byte) (err error) {
	err = t.PutInverted([]byte(token), values)
	return errors.WithMessage(err, "putInverted error")
}

// PutInverted 插入term
func (t *InvertedDB) PutInverted(key, value []byte) (err error) {
	err = Put(t.db, consts.InvertedBucket, key, value)
	return errors.WithMessage(err, "put error")
}

// GetInverted 通过term获取value
func (t *InvertedDB) GetInverted(key []byte) (value []byte, err error) {
	value, err = Get(t.db, consts.InvertedBucket, key)
	err = errors.WithMessage(err, "get error")
	return
}

func (t *InvertedDB) GetAllInverted() (p []*types.InvertedInfo, err error) {
	return
}

// GetInvertedInfo 获取倒排地址
func (t *InvertedDB) GetInvertedInfo(token string) (p *types.InvertedInfo, err error) {
	c, err := t.GetInverted([]byte(token))
	if err != nil {
		err = errors.WithMessage(err, "getInverted error")
		return
	}

	if len(c) == 0 {
		err = errors.Wrap(errors.New("暂无此token"), "len(c) equal to zero")
		return
	}
	output := roaring.New()
	err = output.UnmarshalBinary(c)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshalBinary")
	}
	p = &types.InvertedInfo{
		Token:  token,
		DocIds: output,
	}
	return
}

// GetInvertedDoc 根据地址获取读取文件
func (t *InvertedDB) GetInvertedDoc(offset int64, size int64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), offset/int64(page), int(offset+size))
	if err != nil {
		return nil, errors.WithMessage(errors.Errorf("GetDocinfo Mmap err: %v", err), "mmap error")
	}
	return b[offset : offset+size], nil
}

func (t *InvertedDB) Close() {
	err := t.file.Close()
	if err != nil {
		return
	}
	err = t.db.Close()
	if err != nil {
		return
	}
}
