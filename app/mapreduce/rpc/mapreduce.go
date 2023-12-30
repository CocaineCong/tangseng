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
