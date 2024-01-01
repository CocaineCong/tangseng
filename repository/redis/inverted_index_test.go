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

package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	InitRedis()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestPushInvertedPath(t *testing.T) {
	ctx := context.Background()
	_ = PushInvertedPath(ctx, InvertedIndexDbPathDayKey, []string{"a", "b", "c"})
}

func TestListInvertedPath(t *testing.T) {
	ctx := context.Background()
	paths, _ := ListInvertedPath(ctx, []string{InvertedIndexDbPathDayKey})
	fmt.Println(paths)
}

func TestSetInvertedIndexTokenDocIds(t *testing.T) {
	ctx := context.Background()
	docIds := roaring.NewBitmap()
	docIds.AddInt(1)
	docIds.AddInt(2)
	err := SetInvertedIndexTokenDocIds(ctx, "test", docIds)
	fmt.Println(err)
}

func TestGetInvertedIndexTokenDocIds(t *testing.T) {
	ctx := context.Background()
	docIds, err := GetInvertedIndexTokenDocIds(ctx, "test1")
	fmt.Println(err)
	fmt.Println(docIds)
}

func TestPushInvertedIndexToken(t *testing.T) {
	ctx := context.Background()
	err := PushInvertedIndexToken(ctx, 1, "test2")
	fmt.Println(err)
}

func TestGetInvertedIndexToken(t *testing.T) {
	ctx := context.Background()
	tokens, err := ListInvertedIndexToken(ctx, 1)
	fmt.Println(err)
	fmt.Println(tokens)
}

func TestSetInvertedIndexByKey(t *testing.T) {
	ctx := context.Background()
	key := GetInvertedIndexDbPathMonthKey("10")
	_ = SetInvertedPath(ctx, key, "a1")
	key2 := GetInvertedIndexDbPathMonthKey("11")
	_ = SetInvertedPath(ctx, key2, "b")
}

func TestListInvertedIndexByPrefixKey(t *testing.T) {
	ctx := context.Background()
	key := GetInvertedIndexDbPathMonthKey("*")
	result, _ := ListInvertedIndexByPrefixKey(ctx, key)
	fmt.Println(result)
}
