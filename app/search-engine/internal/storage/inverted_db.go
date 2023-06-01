package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/boltdb/bolt"

	log "github.com/CocaineCong/Go-SearchEngine/pkg/logger"
	"github.com/CocaineCong/Go-SearchEngine/pkg/util/se"
)

const termBucket = "term"

type KvInfo struct {
	Key   []byte
	Value []byte
}

type InvertedDB struct {
	file   *os.File
	db     *bolt.DB
	offset int64
}

type TermValue struct {
	DocCount int64
	Offset   int64
	Size     int64
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
	log.LogrusObj.Infof("start op bolt:%s", termName)
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		log.LogrusObj.Error(err)
	}
	return &InvertedDB{f, db, stat.Size()}
}

func (t *InvertedDB) StoragePostings(token string, values []byte, docCount int64) (err error) {
	// 写入file，获取写入的size
	size, err := t.storagePostings(values)
	if err != nil {
		return
	}
	// 写入b+ tree
	buf := bytes.NewBuffer(nil)
	err = se.BinaryWrite(buf, docCount)
	if err != nil {
		return
	}
	err = se.BinaryWrite(buf, []int64{t.offset, size})
	if err != nil {
		return
	}

	t.offset += size
	return t.PutInverted([]byte(token), buf.Bytes())
}

// PutInverted 插入term
func (t *InvertedDB) PutInverted(key, value []byte) error {
	return Put(t.db, termBucket, key, value)
}

// GetInverted 通过term获取value
func (t *InvertedDB) GetInverted(key []byte) (value []byte, err error) {
	return Get(t.db, termBucket, key)
}

// GetTermInfo 获取term关联的倒排地址
func (t *InvertedDB) GetTermInfo(token string) (*TermValue, error) {
	c, err := t.GetInverted([]byte(token))
	if err != nil {
		return nil, err
	}
	return Bytes2TermVal(c)
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

// GetInvertedTermCursor 获取遍历游标
func (t *InvertedDB) GetInvertedTermCursor(ternCH chan KvInfo) error {
	return t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(termBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ternCH <- KvInfo{k, v}
		}
		close(ternCH)
		return nil
	})
}

func (t *InvertedDB) storagePostings(postings []byte) (size int64, err error) {
	s, err := t.file.WriteAt(postings, t.offset)
	if err != nil {
		return 0, err
	}

	return int64(s), nil
}

func (t *InvertedDB) Close() {
	t.file.Close()
	t.db.Close()
}

// Bytes2TermVal 字节转换为TermValues
func Bytes2TermVal(values []byte) (p *TermValue, err error) {
	err = binary.Read(bytes.NewBuffer(values), binary.LittleEndian, p)
	if err != nil {
		return
	}

	return
}
