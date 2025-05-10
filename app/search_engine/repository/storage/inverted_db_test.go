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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/CocaineCong/tangseng/config"
)

func TestGetInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	p, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}

// 为redisMockChan产生文件名
func getMsg(testDir string, start, end int) []string {
	msg := make([]string, 0)
	for file_id := start; file_id < end; file_id++ {
		msg = append(msg, testDir+fmt.Sprintf("%d", file_id))
	}
	return msg
}

func TestInitInvertedDB(t *testing.T) {
	testDir := "/tmp/ts/TestInvertedDBManager/"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			return
		}
	}(testDir) // 确保在测试结束时删除目录
	mockRedisChan := make(chan []string, 10)
	// TODO: ctx的操作应该都放到context中
	ctx := context.Background()
	ctx = context.WithValue(ctx, "cleanTime", 2)                 //nolint:all
	ctx = context.WithValue(ctx, "mockRedisChan", mockRedisChan) //nolint:all

	// 向channel发送数据
	mockRedisChan <- getMsg(testDir, 3, 10)

	InitInvertedDB(ctx)

	// 睡眠3秒，确保后台clean线程删除version-0
	time.Sleep(3 * time.Second)
	if len(GlobalInvertedDB.versionSet) != 1 {
		t.Errorf("Expected %v, but got %v", 1, len(GlobalInvertedDB.versionSet))
	}

	if GlobalInvertedDB.currentVersion.versionId != 1 {
		t.Errorf("Expected %v, but got %v", 1, GlobalInvertedDB.currentVersion.versionId)
	}
	// 使用当前版本，然后新建一个版本
	_, oldVersionId := GlobalInvertedDB.Ref()
	if GlobalInvertedDB.currentVersion.ref.Load() != 1 {
		t.Errorf("Expected %v, but got %v", 1, GlobalInvertedDB.currentVersion.ref.Load())
	}
	mockRedisChan <- getMsg(testDir, 6, 12)
	GlobalInvertedDB.UpdateFromRedis(ctx)

	/*
	   此时有两个version
	                                   current ↓
	       | version-1,ref:1 | version-2,ref:0 |
	*/
	if GlobalInvertedDB.currentVersion.versionId != 2 {
		t.Errorf("Expected %v, but got %v", 1, GlobalInvertedDB.currentVersion.versionId)
	}

	// 去引用，这个版本会被异步释放掉
	GlobalInvertedDB.Unref(oldVersionId)

	// 睡眠3秒，确保后台clean线程删除version-1
	time.Sleep(3 * time.Second)
	// 当前只有最新版本
	if GlobalInvertedDB.currentVersion.versionId != 2 {
		t.Errorf("Expected %v, but got %v", 1, GlobalInvertedDB.currentVersion.versionId)
	}
	if len(GlobalInvertedDB.versionSet) != 1 {
		t.Errorf("Expected %v, but got %v", 1, len(GlobalInvertedDB.versionSet))
	}
}
