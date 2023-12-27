package master

import (
	"context"
	"github.com/pkg/errors"
	"math"
	"sync"
	"time"

	"github.com/CocaineCong/tangseng/app/mapreduce/consts"
	"github.com/CocaineCong/tangseng/consts/e"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/types"
)

type MasterSrv struct {
	TaskQueue     chan *types.MapReduceTask // 等待执行的task
	TaskMeta      map[int]*types.MasterTask // 当前所有task的信息
	MasterPhase   types.State               // Master的阶段
	NReduce       int                       // Reduce的数量
	InputFiles    []string                  // 输入的文件
	Intermediates [][]string                // Map任务产生的R个中间文件的信息

	mapreduce.UnimplementedMapReduceServiceServer
}

var (
	InputFiles = []string{ // TODO 配置文件读取
		"/Users/mac/GolandProjects/Go-SearchEngine/app/mapreduce/input_data/other_input_data/movies_data.csv",
		"/Users/mac/GolandProjects/Go-SearchEngine/app/mapreduce/input_data/other_input_data/movies_data1.csv",
		"/Users/mac/GolandProjects/Go-SearchEngine/app/mapreduce/input_data/other_input_data/movies_data2.csv",
	}
	MapReduceSrvIns  *MasterSrv
	MapReduceSrvOnce sync.Once
	mu               sync.Mutex
)

func GetMapReduceSrv() *MasterSrv {
	MapReduceSrvOnce.Do(func() {
		MapReduceSrvIns = NewMaster(InputFiles, consts.ReduceDefaultNum)
	})
	return MapReduceSrvIns
}

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

func (m *MasterSrv) createMapTask() {
	for idx, filename := range m.InputFiles {
		// 把输入的files都形成一个task元数据塞到queue中
		taskMeta := types.MapReduceTask{
			Input:      filename,
			TaskState:  types.Map,
			NReducer:   m.NReduce,
			TaskNumber: idx,
		}
		m.TaskQueue <- &taskMeta
		m.TaskMeta[idx] = &types.MasterTask{
			TaskStatus:    types.Idle, // 状态为 idle ，等待worker节点来进行
			TaskReference: &taskMeta,
			StartTime:     time.Now(),
		}
	}
}

func (m *MasterSrv) Done() bool {
	mu.Lock()
	defer mu.Unlock()
	return m.MasterPhase == types.Exit
}

func (m *MasterSrv) catchTimeout() { // nolint:golint,unused
	for {
		time.Sleep(5 * time.Second)
		mu.Lock()
		if m.MasterPhase == types.Exit {
			mu.Unlock()
			return
		}
		for _, masterTask := range m.TaskMeta {
			if masterTask.TaskStatus == types.InProgress &&
				time.Since(masterTask.StartTime) > 10*time.Second {
				m.TaskQueue <- masterTask.TaskReference
				masterTask.TaskStatus = types.Idle
			}
		}
		mu.Unlock()
	}
}

func (m *MasterSrv) createReduceTask() {
	m.TaskMeta = map[int]*types.MasterTask{}
	for idx, files := range m.Intermediates {
		taskMeta := types.MapReduceTask{
			TaskState:     types.Reduce,
			NReducer:      m.NReduce,
			TaskNumber:    idx,
			Intermediates: files,
		}
		m.TaskQueue <- &taskMeta
		m.TaskMeta[idx] = &types.MasterTask{
			TaskStatus:    types.Idle,
			TaskReference: &taskMeta,
		}
	}
}

// MasterAssignTask master 等待worker调用
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
		// 有就发出去
		*task = *<-m.TaskQueue
		// 记录task的启动时间
		m.TaskMeta[task.TaskNumber].TaskStatus = types.InProgress
		m.TaskMeta[task.TaskNumber].StartTime = time.Now()
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
		err = errors.WithMessage(err, "processTaskResult error")
		resp.Code = e.ERROR
		resp.Message = "map finish failed"
		return
	}

	return
}

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

func (m *MasterSrv) allTaskDone() bool {
	for _, task := range m.TaskMeta {
		if task.TaskStatus != types.Completed {
			return false
		}
	}

	return true
}
