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
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

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

var GlobalInvertedDB *InvertedDBManager

type InvertedDB struct {
	file   *os.File
	db     *bolt.DB
	offset int64
}

// 使用manager来管理多个倒排db
// currentVersion为当前最新的版本
// versionSet 记录了还在使用版本的信息
type InvertedDBManager struct {
	currentVersion *Version
	versionSet     map[int64]*Version
	rwMutex        sync.RWMutex
}

// 当前的版本信息
type Version struct {
	versionId int64
	dbTable   map[string]*InvertedDB
	oldFiles  []string
	dbs       []*InvertedDB
	ref       atomic.Int64
}

// InitInvertedDB 初始化倒排索引库
func InitInvertedDB(ctx context.Context) {
	// 新建一个dummy version
	version := Version{
		versionId: 0,
		dbTable:   make(map[string]*InvertedDB),
		dbs:       make([]*InvertedDB, 0),
	}
	manager := InvertedDBManager{
		currentVersion: &version,
		versionSet:     make(map[int64]*Version),
	}
	manager.versionSet[0] = &version
	manager.UpdateFromRedis(ctx)
	// cleanTime主要用来测试用
	cleanTime, ok := ctx.Value("cleanTime").(int)
	if !ok {
		// 未设置就按照30分钟来异步清理一次
		go manager.backgroundCleaner(30 * 60)
	} else {
		go manager.backgroundCleaner(cleanTime)
	}
	GlobalInvertedDB = &manager
}

func (m *InvertedDBManager) UpdateFromRedis(ctx context.Context) {
	newTable := make(map[string]*InvertedDB)
	// mock redis用于测试，测试ok可以丢掉这部分
	mockRedisChan, ok := ctx.Value("mockRedisChan").(chan []string)
	var filePath []string
	if !ok {
		filePath, _ = redis.ListInvertedPath(ctx, redis.InvertedIndexDbPathKeys)
	} else {
		filePath = <-mockRedisChan
	}
	// 使用newTable异地构建，可仅上读锁不影响在线服务
	m.rwMutex.RLock()
	// 新建version
	version := Version{
		versionId: m.currentVersion.versionId + 1,
		dbs:       make([]*InvertedDB, 0),
	}
	version.ref.Store(0)

	// 找到没有被映射的file，构建倒排db
	// 复用还在使用的db
	currentDBTable := m.currentVersion.dbTable
	for _, file := range filePath {
		_, exist := currentDBTable[file]
		if !exist {
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
			currentDBTable[file] = &InvertedDB{f, db, stat.Size()}
		}
		// 添加倒排db指针
		idb := currentDBTable[file]
		newTable[file] = idb
		version.dbs = append(version.dbs, idb)
	}

	// 找到新版本不再使用的file，记录下来
	oldFiles := make([]string, 0)
	for file := range currentDBTable {
		// newTable中没有旧的file
		if _, exist := newTable[file]; !exist {
			oldFiles = append(oldFiles, file)
		}
	}
	// 设置当前版本的dbTable
	version.dbTable = newTable
	m.rwMutex.RUnlock()

	// 上写锁更新版本，这里会影响到在线服务，尽量减少写锁内的操作
	m.rwMutex.Lock()
	// 新版本不再使用的file放到当前版本
	m.versionSet[m.currentVersion.versionId].oldFiles = oldFiles
	m.versionSet[version.versionId] = &version
	m.currentVersion = &version
	m.rwMutex.Unlock()
}

// 后台异步清理掉不再使用的倒排db，构建操作不应该是频繁操作
func (m *InvertedDBManager) backgroundCleaner(cleanTime int) {
	// 接受信号优雅退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		// 用计算器来每cleanTime秒检查一次是否有旧版本需要清理
		timer := time.NewTimer(time.Duration(cleanTime) * time.Second)
		defer timer.Stop()
		select {
		case <-sig:
			return
		case <-timer.C:
			m.cleanOldVersion()
		}
	}
}

func (m *InvertedDBManager)cleanOldVersion() {
	oldIds := make([]int64, 0)
	m.rwMutex.RLock()
	if len(m.versionSet) > 1 {
		//存在旧版本，关闭不再需要的db
		for id, version := range m.versionSet {
			// 为什么可以直接清理掉ref为0的版本(最新版本除外)?
			// 因为旧版本不可能再被用到,后续请求永远使用最新版本
			if id < m.currentVersion.versionId && version.ref.Load() == 0 {
				for _, file := range version.oldFiles {
					if idb, exist := version.dbTable[file]; exist {
						idb.Close()
					}
				}
				//记录一下，后续在版本链中清除
				oldIds = append(oldIds, id)
			}
		}
	}
	m.rwMutex.RUnlock()
	//版本链中去除旧版本需要写锁
	m.rwMutex.Lock()
	for _, oldId := range oldIds {
		delete(m.versionSet, oldId)
	}
	m.rwMutex.Unlock()
}

// 仅能通过ref来获取倒排db
func (m *InvertedDBManager) Ref() ([]*InvertedDB, int64) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	//获取当前版本files并添加引用
	versionId := m.currentVersion.versionId
	dbs := m.currentVersion.dbs
	m.currentVersion.ref.Add(1)
	return dbs, versionId
}

// 用完需要回收
func (m *InvertedDBManager) Unref(versionId int64) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	version := m.versionSet[versionId]
	version.ref.Add(-1)
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
