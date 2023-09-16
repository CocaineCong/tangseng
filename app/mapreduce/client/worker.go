package client

import (
	"context"
	"hash/fnv"
	"time"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/app/mapreduce/rpc"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

func Worker(ctx context.Context, mapf func(string, string) []*types.KeyValue, reducef func(string, []string) *roaring.Bitmap) {
	// 启动worker
	for {
		// worker从master获取任务
		task, err := getTask(ctx)
		if err != nil {
			log.LogrusObj.Error("Worker-getTask", err)
			return
		}
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
	h.Write([]byte(key))
	return int64(h.Sum32() & 0x7fffffff)
}

func getTask(ctx context.Context) (resp *mapreduce.MapReduceTask, err error) {
	// worker从master获取任务
	taskReq := &mapreduce.MapReduceTask{}
	resp, err = rpc.MapReduceClient.MasterAssignTask(ctx, taskReq)

	return
}

func TaskCompleted(ctx context.Context, task *mapreduce.MapReduceTask) (reply *mapreduce.MasterTaskCompletedResp, err error) {
	// 通过RPC，把task信息发给master
	reply, err = rpc.MapReduceClient.MasterTaskCompleted(ctx, task)

	return
}
