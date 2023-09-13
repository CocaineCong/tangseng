package mapreduce

import (
	"hash/fnv"
	"time"

	"github.com/CocaineCong/tangseng/types"
)

func Worker(mapf func(string, string) []types.KeyValue, reducef func(string, []string) string) {
	// 启动worker
	for {
		// worker从master获取任务
		task := getTask()

		// 拿到task之后，根据task的state，map task交给mapper， reduce task交给reducer
		// 额外加两个state，让 worker 等待 或者 直接退出
		switch task.TaskState {
		case types.Map:
			mapper(&task, mapf)
		case types.Reduce:
			reducer(&task, reducef)
		case types.Wait:
			time.Sleep(5 * time.Second)
		case types.Exit:
			return
		}
	}
}

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func getTask() types.MapReduceTask {
	// worker从master获取任务
	args := types.ExampleArgs{}
	reply := types.MapReduceTask{}
	call("Master.AssignTask", &args, &reply)
	return reply
}

func TaskCompleted(task *types.MapReduceTask) {
	// 通过RPC，把task信息发给master
	reply := types.ExampleReply{}
	call("Master.TaskCompleted", task, &reply)
}
