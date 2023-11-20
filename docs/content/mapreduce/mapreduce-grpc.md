# MapReduce

## 写在前面

**原论文地址：[MapReduce: Simplified Data Processing on Large Clusters](https://research.google.com/archive/mapreduce-osdi04.pdf)**

## 总览

这次 lab1 的 mapreduce，其实是在写 [搜索引擎tangseng][https://github.com/CocaineCong/tangseng] 的时候，需要用来构建倒排索引。所以会和课程上所要求的不太一样，这里也没有使用rpc调用，而是为了与项目统一，便改用了grpc进行调用。

![mapreduce工作原理](https://raw.githubusercontent.com/CremeU/cloud-img/main/image-20230919010122185.png)

这里需要注意几点

- 不同的Map任务之间不会进行通信
- 不同的Reduce任务之间也不会发生任何信息交换
- 所有的数据交换都是通过MapReduce框架自身去实现的

那么如何对 map tasks 和 reduce tasks 进行合理的协调呢？这里我们就要引入两个角色，**master 和 worker**，在原论文中，对这两者的并没有非常明确的定义，但我们可以摘录并提炼原论文对这两个角色的描述：

**master : The master picks idle workers and assigns each one a map task or a reduce task.**

**worker :**

- **The map worker who is assigned a map task reads the contents of the corresponding input split.**

- **The reduce worker iterates over the sorted intermediate data and for each unique intermediate key encountered, it passes the key and the corresponding set of intermediate values to the user’s Reduce function**

> 这里我们先说一下几个状态枚举值：
>
> - idle ：空闲状态
> - in-progress ：进行状态
> - completed ：完成状态
>
> 这三个枚举值代表着每一个 map task 和 reduce task 的状态，标识着这些 task 是未开始，进行中，还是已完成。

那么 master 其实就是选择空闲的 worker 节点，为每一个空闲的 worker 节点分配 map task 或者 reduce task。而 worker 看似分成了 map worker 和 reduce worker，但其实这两个 worker 都是一样，只是看 master 分配的是 map task 还是 reduce task。这样我们的 map 和 reduce 的数据传送就非常清晰了。

![MapReduce整体工作流](https://raw.githubusercontent.com/CremeU/cloud-img/main/image-20230919014314364.png)

接下来，我们来详细讲解一下这几个重要的角色

## Worker

首先我们先定义一个 MapReduce 的任务，也就是我们 worker 需要用到参数

```go
type MapReduceTask struct {
	Input         string   `json:"input"`         // 输入的文件
	TaskState     State    `json:"task_state"`    // 状态
	NReducer      int      `json:"n_reducer"`     // reducer 数量
	TaskNumber    int      `json:"task_number"`   // 任务数量
	Intermediates []string `json:"intermediates"` // map 之后的文件存储地址
	Output        string   `json:"output"`        // output的输出地址
}
```

接着再定义 State 枚举值

```go
type MasterTaskStatus int

const (
	Idle       MasterTaskStatus = iota + 1 // 未开始
	InProgress                             // 进行中
	Completed                              // 已完成
)
```

接下来我们的 Worker 函数就很简单了

```go
func Worker(ctx context.Context, mapf func(string, string) []*types.KeyValue, reducef func(string, []string) *roaring.Bitmap) {
	// 启动worker
  for {
		task, err := getTask(ctx) // worker从master获取任务
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
```

至于 mapper 和 reducer 如何实现的，先桥豆麻袋一下，下文在 map 和 reduce 中会给出答案，如何从 master 中拿到 task 呢？这就涉及到 worker 和 master 的通信。本来打算用 RPC 通信的，但为了项目的整体统一，还是用了 gRPC 。

创建一个proto文件

```protobuf
syntax="proto3";
option go_package = "/index_platform;";

message MapReduceTask{
	// @inject_tag:form:"input" uri:"input"
	string input = 1;
	// @inject_tag:form:"task_state" uri:"task_state"
	int64 task_state = 2;
	// @inject_tag:form:"n_reducer" uri:"n_reducer"
	int64 n_reducer = 3;
	// @inject_tag:form:"task_number" uri:"task_number"
	int64 task_number = 4;
	// @inject_tag:form:"intermediates" uri:"intermediates"
	repeated string intermediates = 5;
	// @inject_tag:form:"output" uri:"output"
	string output = 6;
}

message MasterTaskCompletedResp {
	// @inject_tag:form:"code" uri:"code"
	int64 code=1;
	// @inject_tag:form:"message" uri:"message"
	string message=2;
}

service MapReduceService {
	rpc MasterAssignTask(MapReduceTask) returns (MapReduceTask);
	rpc MasterTaskCompleted(MapReduceTask) returns (MasterTaskCompletedResp);
}
```

定义两个 RPC 函数，`MasterAssignTask `  用来接受 master 分配的 task   `MasterTaskCompleted` 完成 task 之后，对这个 task 进行标识，意味着该任务结束。

所以我们 worker 接受任务的通信如下

```go
func getTask(ctx context.Context) (resp *mapreduce.MapReduceTask, err error) {
	// worker从master获取任务
	taskReq := &mapreduce.MapReduceTask{}
	resp, err = rpc.MapReduceClient.MasterAssignTask(ctx, taskReq)

	return
}
```

当完成任务时，通过gRPC发送给master

```go
func TaskCompleted(ctx context.Context, task *mapreduce.MapReduceTask) (reply *mapreduce.MasterTaskCompletedResp, err error) {
	// 通过RPC，把task信息发给master
	reply, err = rpc.MapReduceClient.MasterTaskCompleted(ctx, task)

	return
}
```

那么 master 是如何分配任务的？接下来我们来介绍一下 master 节点。

## Master

我们定义这么一个 Master 服务的结构体

```go
type MasterSrv struct {
	TaskQueue     chan *types.MapReduceTask // 等待执行的task
	TaskMeta      map[int]*types.MasterTask // 当前所有task的信息
	MasterPhase   types.State               // Master的阶段
	NReduce       int                       // Reduce的数量
	InputFiles    []string                  // 输入的文件
	Intermediates [][]string                // Map任务产生的R个中间文件的信息

	mapreduce.UnimplementedMapReduceServiceServer // gRPC服务实现接口
}
```

那么当我们 New 一个 Master 服务的时候，顺便创建 map tasks 任务

```go
func NewMaster(files []string, nReduce int) *MasterSrv {
	m := &MasterSrv{
		TaskQueue:     make(chan *types.MapReduceTask, int(math.Max(float64(nReduce), float64(len(files))))),
		TaskMeta:      map[int]*types.MasterTask{},
		MasterPhase:   types.Map,
		NReduce:       nReduce,
		InputFiles:    files,
		Intermediates: make([][]string, nReduce),
	}
	m.createMapTask()
	return m
}
```

创建 map task 任务

```go
func (m *MasterSrv) createMapTask() {
  // 把输入的files都形成一个task元数据塞到queue中
	for idx, filename := range m.InputFiles { 
		taskMeta := types.MapReduceTask{
			Input:      filename,
			TaskState:  types.Map, // map节点
			NReducer:   m.NReduce,
			TaskNumber: idx,
		}
		m.TaskQueue <- &taskMeta
		m.TaskMeta[idx] = &types.MasterTask{
			TaskStatus:    types.Idle, // 状态为 idle ，等待worker节点来领取 task
			TaskReference: &taskMeta,
		}
	}
}
```

创建 reduce task 任务

```go
func (m *MasterSrv) createReduceTask() {
	m.TaskMeta = map[int]*types.MasterTask{}
	for idx, files := range m.Intermediates {
		taskMeta := types.MapReduceTask{
			TaskState:     types.Reduce, // reduce 阶段
			NReducer:      m.NReduce,
			TaskNumber:    idx,
			Intermediates: files,
		}
		m.TaskQueue <- &taskMeta
		m.TaskMeta[idx] = &types.MasterTask{
			TaskStatus:    types.Idle, // 找到空闲的 worker
			TaskReference: &taskMeta,
		}
	}
}
```

MasterAssignTask 等待 worker 来领取 task

```go
func (m *MasterSrv) MasterAssignTask(ctx context.Context, req *mapreduce.MapReduceTask) (reply *mapreduce.MapReduceTask, err error) {
	mu.Lock()
	defer mu.Unlock()
	task := &types.MapReduceTask{
		Input:         req.Input,
		TaskState:     types.State(req.TaskState),
		NReducer:      int(req.NReducer),
		TaskNumber:    int(req.TaskNumber),
		Intermediates: req.Intermediates,
		Output:        req.Output,
	}
	if len(m.TaskQueue) > 0 {
		// 如果queue中还有任务的话就发出去
		*task = *<-m.TaskQueue
		m.TaskMeta[task.TaskNumber].TaskStatus = types.InProgress // 修改worker的状态为进行中
		m.TaskMeta[task.TaskNumber].StartTime = time.Now() // 记录task的启动时间
	} else if m.MasterPhase == types.Exit {
		*task = types.MapReduceTask{
			TaskState: types.Exit,
		}
	} else {
		// 没有task就让worker等待
		*task = types.MapReduceTask{TaskState: types.Wait}
	}
	
  // 返回该任务的状态，因为发出去就是给task了，这个状态已经改变了，worker可以工作了
	reply = &mapreduce.MapReduceTask{
		Input:         task.Input,
		TaskState:     int64(task.TaskState),
		NReducer:      int64(task.NReducer),
		TaskNumber:    int64(task.TaskNumber),
		Intermediates: task.Intermediates,
		Output:        task.Output,
	}

	return
}
```

那么如果 task 把任务都做完了，master 应该怎么回应呢？

```go
func (m *MasterSrv) MasterTaskCompleted(ctx context.Context, req *mapreduce.MapReduceTask) (resp *mapreduce.MasterTaskCompletedResp, err error) {
	resp = new(mapreduce.MasterTaskCompletedResp)
	resp.Code = e.ERROR
	resp.Message = "map finish successfully"
	// 更新task状态
	if req.TaskState != int64(m.MasterPhase) || m.TaskMeta[int(req.TaskNumber)].TaskStatus == types.Completed {
		// 因为worker写在同一个文件这次盘上对于重复的结果要丢弃
		return
	}
	m.TaskMeta[int(req.TaskNumber)].TaskStatus = types.Completed
	err = m.processTaskResult(req) // always success haha and hope u so :)
	if err != nil {
		resp.Code = e.ERROR
		resp.Message = "map finish failed"
		return
	}

	return
}
```

处理任务的结果，如果是 map 完成后就变成 reduce 阶段，reduce 之后就是 all done.  :)

```go
// processTaskResult 处理任务结果
func (m *MasterSrv) processTaskResult(task *mapreduce.MapReduceTask) (err error) {
	switch task.TaskState {
	case int64(types.Map):
		// 收集intermediate信息
		for reduceTaskId, filePath := range task.Intermediates {
			m.Intermediates[reduceTaskId] = append(m.Intermediates[reduceTaskId], filePath)
		}
		if m.allTaskDone() {
			// 获取所有的map task后，进入reduce阶段
			m.createReduceTask()
			m.MasterPhase = types.Reduce
		}
	case int64(types.Reduce):
		if m.allTaskDone() {
			// 获得所有的reduce task后，进去exit阶段
			m.MasterPhase = types.Exit
		}
	}

	return
}
```

介绍完master之后，我们具体来看一下map的具体行为。

## Map

在 map 中，我们抽离出一个 mapper，具体的map函数可根据实际情况进行修改，然后将map function传入mapper中进行实际的map动作，我们读取每一个文件，然后把输出的结果都放到 `intermediates` 中，并且根据 task 所设定的 `NReducer` 也就是 reducer 数 进行hash ，将结果均匀分到每个中间文件中。

```go
func mapper(ctx context.Context, task *mapreduce.MapReduceTask, mapf func(string, string) []*types.KeyValue) {
	// 从文件名读取content
	content, err := os.ReadFile(task.Input)
	if err != nil {
		log.LogrusObj.Error("mapper", err)
		return
	}
	// 将content交给mapf，缓存结果
	intermediates := mapf(task.Input, string(content))

	// 缓存后的结果会写到本地磁盘，并切成R份
	// 切分方式是根据key做hash
	buffer := make([][]*types.KeyValue, task.NReducer)
	for _, intermediate := range intermediates {
		slot := ihash(intermediate.Key) % task.NReducer
		buffer[slot] = append(buffer[slot], intermediate)
	}
	mapOutput := make([]string, 0)
	for i := 0; i < int(task.NReducer); i++ {
		mapOutput = append(mapOutput, writeToLocalFile(int(task.TaskNumber), i, &buffer[i]))
	}
	// R个文件的位置发送给master
	task.Intermediates = mapOutput
	_, err = TaskCompleted(ctx, task) // 完成后，给master发送消息，map阶段结束
	if err != nil {
		fmt.Println("mapper-TaskCompleted", err)
	}

	return
}
```

具体的 Map方法，由于是用于搜索引擎，所以这里是建立倒排索引

```go
func Map(filename string, contents string) (res []*types.KeyValue) {
	res = make([]*types.KeyValue, 0)
	lines := strings.Split(contents, "\r\n") // 分行
	var inputData *model.InputData
	for _, line := range lines[1:] {
		docStruct, _ := doc2Struct(line) // 字符串转 doc struct
		tokens, err := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
		if err != nil {
			return
		}
		for _, v := range tokens {
      res = append(res, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)}) // token:docId 倒排索引
		}
	}

	return
}
```

至此map就已经完成了，是不是很简单，其实具体的map和reduce并不难，难的是如何平衡调度，接下来我们来看看reduce是如何怎么的。

## Reduce

和map一样，我们抽离出一个reducer，然后把具体的 reduce 传进去，当然还有一个shuffle过程，这里进行排序会减少后面的reduce计算。可以少计算几次。

```go
func reducer(ctx context.Context, task *mapreduce.MapReduceTask, reducef func(string, []string) *roaring.Bitmap) {
	// 先从filepath读取intermediate的KeyValue
	intermediate := *readFromLocalFile(task.Intermediates)
	// 根据kv排序 shuffle 过程
	sort.Sort(types.ByKey(intermediate))

	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/mr-tmp-%d.%s",
		dir, task.TaskNumber, consts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	output := roaring.NewBitmap()
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
		output = reducef(intermediate[i].Key, values)

		// 落倒排索引库
		outByte, _ = output.MarshalBinary()
		_ = invertedDB.StoragePostings(intermediate[i].Key, outByte)
		i = j
	}

	task.Output = outName
	_, err := TaskCompleted(ctx, task) // 完成后，给master发送消息，reduce阶段结束
	if err != nil {
		fmt.Println("reducer-TaskCompleted", err)
		return
	}
}
```

具体的Reduce，其实就是把相同的key的value聚合在一起。比如

after map:

```shell
{"apple":1}
{"apple:"2}
{"poizon":3}
```

after reduce:

```she
{"apple":{1,2}}
{"poizon":{3}}
```

具体实现如下所示：

```go
func Reduce(key string, values []string) *roaring.Bitmap {
	docIds := roaring.New()
	for _, v := range values {
		docIds.AddInt(cast.ToInt(v))
	}
	return docIds
}
```

最终 output 输出

![output](https://raw.githubusercontent.com/CremeU/cloud-img/main/1695145441347.jpg)
