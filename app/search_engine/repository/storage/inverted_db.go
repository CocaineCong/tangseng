package storage

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/fileutils"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

type KvInfo struct {
	Key   []byte
	Value []byte
}

var GobalInvertedDB []*InvertedDB

type InvertedDB struct {
	file   *os.File
	db     *bolt.DB
	offset int64
}

// const InvertedDBPaths = "../../../index_platform/index_data/"
const InvertedDBPaths = "/Users/mac/GolandProjects/Go-SearchEngine/app/index_platform/index_data/"

// InitInvertedDB 初始化倒排索引库
func InitInvertedDB() []*InvertedDB {
	dbs := make([]*InvertedDB, 0)
	filePath := fileutils.GetFiles(InvertedDBPaths)
	for _, file := range filePath {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.LogrusObj.Error(err)
		}
		stat, err := f.Stat()
		if err != nil {
			log.LogrusObj.Error(err)
		}
		db, err := bolt.Open(file, 0600, nil)
		if err != nil {
			log.LogrusObj.Error(err)
		}
		dbs = append(dbs, &InvertedDB{f, db, stat.Size()})
	}
	GobalInvertedDB = dbs
	return nil
}

// NewInvertedDB 新建一个inverted
func NewInvertedDB(termName, postingsName string) *InvertedDB {
	f, err := os.OpenFile(postingsName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.LogrusObj.Error(err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.LogrusObj.Error(err)
	}
	log.LogrusObj.Infof("start op bolt:%v", termName)
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.LogrusObj.Error(err)
	}
	return &InvertedDB{f, db, stat.Size()}
}

// GetInverted 通过term获取value
func (t *InvertedDB) GetInverted(key []byte) (value []byte, err error) {
	return Get(t.db, consts.InvertedBucket, key)
}

// GetInvertedDoc 根据地址获取读取文件
func (t *InvertedDB) GetInvertedDoc(offset int64, size int64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), offset/int64(page), int(offset+size))
	if err != nil {
		return nil, fmt.Errorf("GetDocinfo Mmap err: %v", err)
	}
	return b[offset : offset+size], nil
}

func (t *InvertedDB) Close() {
	t.file.Close()
	t.db.Close()
}
