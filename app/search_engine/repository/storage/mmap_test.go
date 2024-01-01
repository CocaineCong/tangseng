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

//go:build !windows
// +build !windows

package storage

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestMmap(t *testing.T) {
	filePath := "../../data/db/0.forward"
	fd, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer fd.Close()

	// 获取文件大小
	fi, err := fd.Stat()
	if err != nil {
		fmt.Println("获取文件信息失败:", err)
		return
	}
	fileSize := fi.Size()
	fmt.Println("fileSize", fileSize)
	// 设置要映射的偏移量，假设从中间开始映射
	offset := int64(10)

	// 获取从偏移量开始到结尾的长度
	length := int(fileSize - offset)

	// 映射整个文件到内存
	mmapData, err := Mmap(int(fd.Fd()), offset, length)
	if err != nil {
		fmt.Println("映射文件失败:", err)
		return
	}
	defer func(b []byte) {
		err = syscall.Munmap(b)
		if err != nil {
			fmt.Println(err)
		}
	}(mmapData)

	// 使用 mmapData 可以直接读取文件内容
	fmt.Printf("文件内容：%s\n", string(mmapData))
}
