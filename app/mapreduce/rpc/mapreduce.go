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

package rpc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
)

// MasterAssignTask 通过 master 发送任务
func MasterAssignTask(ctx context.Context, taskReq *mapreduce.MapReduceTask) (resp *mapreduce.MapReduceTask, err error) {
	resp, err = MapReduceClient.MasterAssignTask(ctx, taskReq)
	if err != nil {
		err = errors.WithMessage(err, "MasterAssignTask-MapReduceClient error")
		return
	}

	return
}

// MasterTaskCompleted 通知 master 任务完成的RPC调用
func MasterTaskCompleted(ctx context.Context, task *mapreduce.MapReduceTask) (resp *mapreduce.MasterTaskCompletedResp, err error) {
	resp, err = MapReduceClient.MasterTaskCompleted(ctx, task)
	if err != nil {
		err = errors.WithMessage(err, "MapReduceClient.MasterTaskCompleted error")
		return
	}

	return
}
