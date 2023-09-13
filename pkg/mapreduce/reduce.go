package mapreduce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/CocaineCong/tangseng/types"
)

func reducer(task *types.MapReduceTask, reducef func(string, []string) string) {
	// 先从filepath读取intermediate的KeyValue
	intermediate := *readFromLocalFile(task.Intermediates)
	// 根据kv排序
	sort.Sort(types.ByKey(intermediate))

	dir, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(dir, "mr-tmp-*")
	if err != nil {
		fmt.Println(err)
	}
	// 这部分代码修改自mrsequential.go
	i := 0
	for i < len(intermediate) {
		// 将相同的key放在一起分组合并
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		// 交给reducef，拿到结果
		output := reducef(intermediate[i].Key, values)
		// 写到对应的output文件
		fmt.Fprintf(tempFile, "%v %v\n", intermediate[i].Key, output)
		i = j
	}
	tempFile.Close()
	oname := fmt.Sprintf("mr-out-%d", task.TaskNumber)
	os.Rename(tempFile.Name(), oname)
	task.Output = oname
	TaskCompleted(task)
}

func readFromLocalFile(files []string) *[]types.KeyValue {
	kva := []types.KeyValue{}
	for _, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
		}
		dec := json.NewDecoder(file)
		for {
			var kv types.KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			kva = append(kva, kv)
		}
		file.Close()
	}
	return &kva
}
