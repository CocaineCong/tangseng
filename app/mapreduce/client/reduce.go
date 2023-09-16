package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/types"
)

func reducer(ctx context.Context, task *mapreduce.MapReduceTask, reducef func(string, []string) *roaring.Bitmap) {
	// 先从filepath读取intermediate的KeyValue
	intermediate := *readFromLocalFile(task.Intermediates)
	// 根据kv排序
	sort.Sort(types.ByKey(intermediate))

	dir, _ := os.Getwd()
	tempFile, err := os.CreateTemp(dir, "mr-tmp-*")
	if err != nil {
		fmt.Println(err)
	}
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
		// 写到对应的output文件
		_, _ = fmt.Fprintf(tempFile, "%v %v\n", intermediate[i].Key, output)
		i = j
	}
	_ = tempFile.Close()
	oname := fmt.Sprintf("mr-out-%d", task.TaskNumber)
	_ = os.Rename(tempFile.Name(), oname)
	task.Output = oname
	_, err = TaskCompleted(ctx, task)
	if err != nil {
		fmt.Println("reducer-TaskCompleted", err)
		return
	}
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
		_ = file.Close()
	}
	return &kva
}
