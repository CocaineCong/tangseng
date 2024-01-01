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

package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/types"
)

func reducer(ctx context.Context, task *mapreduce.MapReduceTask, reducef func(string, []string) *roaring.Bitmap) {
	// 先从filepath读取intermediate的KeyValue
	intermediate := *readFromLocalFile(task.Intermediates)
	// 根据kv排序 shuffle 过程
	sort.Sort(types.ByKey(intermediate))

	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/mr-tmp-%d.%s",
		dir, task.TaskNumber, consts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	var outByte []byte

	i := 0
	for i < len(intermediate) {
		// 将相同的key放在一起分组合并
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		var values []string
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		// 交给reducef，拿到结果
		output := reducef(intermediate[i].Key, values)
		// 落倒排索引库
		outByte, _ = output.MarshalBinary()
		_ = invertedDB.StoragePostings(intermediate[i].Key, outByte)
		i = j
	}

	task.Output = outName
	_, err := TaskCompleted(ctx, task)
	if err != nil {
		fmt.Println("reducer-TaskCompleted", err)
		return
	}
}

func readFromLocalFile(files []string) *[]*types.KeyValue {
	kva := []*types.KeyValue{}
	for _, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
		}
		dec := json.NewDecoder(file)
		for {
			var kv *types.KeyValue
			if err = dec.Decode(&kv); err != nil {
				break
			}
			kva = append(kva, kv)
		}
		_ = file.Close()
	}
	return &kva
}
