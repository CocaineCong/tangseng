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

package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/RoaringBitmap/roaring"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/input_data"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	cconsts "github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/consts/e"
	pb "github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/pkg/clone"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/mapreduce"
	"github.com/CocaineCong/tangseng/pkg/timeutils"
	"github.com/CocaineCong/tangseng/pkg/trie"
	"github.com/CocaineCong/tangseng/repository/redis"
	"github.com/CocaineCong/tangseng/types"
)

type IndexPlatformSrv struct {
	*pb.UnimplementedIndexPlatformServiceServer
}

var (
	IndexPlatIns  *IndexPlatformSrv
	IndexPlatOnce sync.Once
)

func GetIndexPlatformSrv() *IndexPlatformSrv {
	IndexPlatOnce.Do(func() {
		IndexPlatIns = new(IndexPlatformSrv)
	})
	return IndexPlatIns
}

// BuildIndexService 构建索引
func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	// 时间估计
	resp = new(pb.BuildIndexResp)
	resp.Code = e.SUCCESS
	resp.Message = e.GetMsg(e.SUCCESS)
	invertedIndex := cmap.New[*roaring.Bitmap]() // 倒排索引
	dictTrie := trie.NewTrie()                   // 前缀树

	logs.LogrusObj.Infof("BuildIndexService Start req: %v", req.FilePath)
	// mapreduce 这个是用chan和goroutine来代替master和worker的rpc调用，避免了频繁的rpc调用
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		for _, path := range req.FilePath {
			content, _ := os.ReadFile(path)
			source <- content
		}
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
		// 控制并发
		var wg sync.WaitGroup
		ch := make(chan struct{}, 3)

		keyValueList := make([]*types.KeyValue, 0, 1e3)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines[1:] {
			ch <- struct{}{}
			wg.Add(1)
			docStruct, _ := input_data.Doc2Struct(line) // line 转 docs struct
			if docStruct.DocId == 0 {
				continue
			}

			// 分词
			tokens, _ := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
			for _, v := range tokens {
				if v.Token == "" || v.Token == " " {
					continue
				}
				keyValueList = append(keyValueList, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
				dictTrie.Insert(v.Token)
			}

			// 建立正排索引
			go func(docStruct *types.Document) {
				err = input_data.DocData2Kfk(docStruct)
				if err != nil {
					logs.LogrusObj.Error(err)
				}
				defer wg.Done()
				<-ch
			}(docStruct)
		}
		wg.Wait()

		// // 构建前缀树 // TODO: kafka异步处理一下前缀树的插入，不然占着这里的资源
		// go func(tokenList []string) {
		// 	err = input_data.DocTrie2Kfk(tokenList)
		// 	if err != nil {
		// 		logs.LogrusObj.Error("DocTrie2Kfk", err)
		// 	}
		// }(tokenList)

		// shuffle 排序过程
		sort.Sort(types.ByKey(keyValueList))
		writer.Write(keyValueList)
	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		for values := range pipe {
			for _, v := range values { // 构建倒排索引
				if value, ok := invertedIndex.Get(v.Key); ok {
					value.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, value)
				} else {
					docIds := roaring.NewBitmap()
					docIds.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, docIds)
				}
			}
		}
	})

	// 存储倒排索引
	go func() {
		newCtx := clone.NewContextWithoutDeadline()
		newCtx.Clone(ctx)
		err = storeInvertedIndexByHash(newCtx, invertedIndex)
		if err != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash error ", err)
		}
	}()

	logs.LogrusObj.Infoln("storeInvertedIndexByHash End")

	// 存储前缀树
	go func() {
		newCtx := clone.NewContextWithoutDeadline()
		newCtx.Clone(ctx)
		err = storeDictTrieByHash(newCtx, dictTrie)
		if err != nil {
			logs.LogrusObj.Error("storeDictTrieByHash error ", err)
			logs.LogrusObj.Errorf("stack trace: \n%+v\n", err)
		}
	}()

	return
}

// storeInvertedIndexByHash 分片存储
func storeInvertedIndexByHash(ctx context.Context, invertedIndex cmap.ConcurrentMap[string, *roaring.Bitmap]) (err error) {
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetNowTime(), cconsts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	// 对所有的key进行存储
	for k, val := range invertedIndex.Items() {
		outByte, errx := val.MarshalBinary()
		if errx != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-MarshalBinary", errx)
			continue
		}
		err = invertedDB.StoragePostings(k, outByte)
		if err != nil {
			logs.LogrusObj.Error("storeInvertedIndexByHash-StoragePostings", err)
			continue
		}
	}
	invertedDB.Close()

	err = redis.PushInvertedPath(ctx, redis.InvertedIndexDbPathDayKey, []string{outName})
	if err != nil {
		return errors.WithMessage(err, "redis.PushInvertedPath error")
	}

	// TODO: hash 分片存储, 目前只是根据天数分库，一天的数据都放到同一个库中，感觉这样还是不太行，还是按照每小时或者ihash进行分库，以下同理
	// dir, _ := os.Getwd()
	// keys := invertedIndex.Keys()
	// buffer := make([][]*types.KeyValue, consts.ReduceDefaultNum)
	// for i, v := range keys {
	// 	val, _ := invertedIndex.Get(v)
	// 	slot := iHash(v) % consts.ReduceDefaultNum
	// 	buffer[slot] = append(buffer[slot])
	// 	fmt.Println(v, val)
	// }
	// outName := fmt.Sprintf("%s/%d.%s", dir, i, cconsts.InvertedBucket)

	return
}

// storeInvertedIndexByHash 分片存储
func storeDictTrieByHash(ctx context.Context, dict *trie.Trie) (err error) {
	// TODO: 抽离一个hash存储的方法
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetNowTime(), cconsts.TrieTreeBucket)
	trieDB := storage.NewTrieDB(outName)
	err = trieDB.StorageDict(dict)
	if err != nil {
		return errors.WithMessage(err, "storageDict error")
	}
	_ = trieDB.Close()

	err = redis.PushInvertedPath(ctx, redis.TireTreeDbPathDayKey, []string{outName})
	if err != nil {
		return errors.WithMessage(err, "redis.PushInvertedPath error")
	}

	return
}

// iHash 哈希作用
func iHash(key string) int64 { // nolint:golint,unused
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int64(h.Sum32() & 0x7fffffff)
}
