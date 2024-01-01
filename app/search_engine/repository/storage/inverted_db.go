// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package storage

import (
	"context"
	"os"

	"github.com/pkg/errors"

	bolt "go.etcd.io/bbolt"

	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/redis"
)

type KvInfo struct {
	Key   []byte
	Value []byte
}

var GlobalInvertedDB []*InvertedDB

type InvertedDB struct {
	file   *os.File
	db     *bolt.DB
	offset int64
}

// InitInvertedDB 初始化倒排索引库
func InitInvertedDB(ctx context.Context) []*InvertedDB {
	dbs := make([]*InvertedDB, 0)
	filePath, _ := redis.ListInvertedPath(ctx, redis.InvertedIndexDbPathKeys)
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
	if len(filePath) == 0 {
		return nil
	}
	GlobalInvertedDB = dbs
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
	value, err = Get(t.db, consts.InvertedBucket, key)
	if err != nil {
		err = errors.WithMessage(err, "get error")
	}
	return
}

// GetInvertedDoc 根据地址获取读取文件
func (t *InvertedDB) GetInvertedDoc(offset int64, size int64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), offset/int64(page), int(offset+size))
	if err != nil {
		return nil, errors.WithMessage(err, "GetDocinfo Mmap error")
	}
	return b[offset : offset+size], nil
}

func (t *InvertedDB) Close() {
	err := t.file.Close()
	if err != nil {
		log.LogrusObj.Error("failed to close file")
	}
	err = t.db.Close()
	if err != nil {
		log.LogrusObj.Error("failed to close db")
	}
}
