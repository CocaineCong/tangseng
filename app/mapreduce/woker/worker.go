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
	"fmt"
	"hash/fnv"
	"time"

	"github.com/pkg/errors"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/app/mapreduce/rpc"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

func Worker(ctx context.Context, mapf func(string, string) []*types.KeyValue, reducef func(string, []string) *roaring.Bitmap) {
	// 启动worker
	fmt.Println("Worker working")
	for {
		// worker从master获取任务
		task, err := getTask(ctx)
		if err != nil {
			log.LogrusObj.Errorf("getTask failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
			log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
			return
		}
		fmt.Println("Worker task", task)
		// 拿到task之后，根据task的state，map task交给mapper， reduce task交给reducer
		// 额外加两个state，让 worker 等待 或者 直接退出
		switch task.TaskState {
		case int64(types.Map):
			mapper(ctx, task, mapf)
		case int64(types.Reduce):
			reducer(ctx, task, reducef)
		case int64(types.Wait):
			time.Sleep(5 * time.Second)
		case int64(types.Exit):
			return
		default:
			return
		}
	}
}

func ihash(key string) int64 {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		log.LogrusObj.Error("failed to write")
		return 0
	}
	return int64(h.Sum32() & 0x7fffffff)
}

func getTask(ctx context.Context) (resp *mapreduce.MapReduceTask, err error) {
	// worker从master获取任务
	fmt.Println("getTask Req")
	taskReq := &mapreduce.MapReduceTask{}
	resp, err = rpc.MasterAssignTask(ctx, taskReq)
	fmt.Println("getTask Resp")
	err = errors.WithMessage(err, "MasterAssignTask error")
	return
}

func TaskCompleted(ctx context.Context, task *mapreduce.MapReduceTask) (reply *mapreduce.MasterTaskCompletedResp, err error) {
	reply, err = rpc.MasterTaskCompleted(ctx, task)
	err = errors.WithMessage(err, "MasterTaskCompleted error")
	return
}
