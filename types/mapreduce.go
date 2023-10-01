package types

import (
	"time"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// for sorting by key.
type ByKey []*KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

type MapReduceTask struct {
	Input         string   `json:"input"`         // 输入的文件
	TaskState     State    `json:"task_state"`    // 状态
	NReducer      int      `json:"n_reducer"`     // reducer 数量
	TaskNumber    int      `json:"task_number"`   // 任务数量
	Intermediates []string `json:"intermediates"` // map 之后的文件存储地址
	Output        string   `json:"output"`        // output的输出地址
}

type MasterTask struct {
	TaskStatus    MasterTaskStatus
	StartTime     time.Time
	TaskReference *MapReduceTask
}

type MasterTaskStatus int

const (
	Idle       MasterTaskStatus = iota + 1 // 未开始
	InProgress                             // 进行中
	Completed                              // 已完成
)

type State int

const (
	Map State = iota + 1
	Reduce
	Exit
	Wait
)

// Tokenization 分词返回结构
type Tokenization struct {
	Token string // 词条
	// Position int64  // 词条在文本的位置 // TODO 后面再补上
	// Offset   int64  // 偏移量
	DocId int64
}
