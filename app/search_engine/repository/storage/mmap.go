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
	"syscall"
)

// Mmap 将一个文件映射到内存中，以便可以直接通过内存访问文件的内容。
// 映射后的内存可以像普通的字节切片一样进行读取和写入操作，而不需要额外的文件读写操作。
// 这对于处理大文件或需要频繁访问文件内容的场景非常有用，因为避免了多次磁盘读写操作，提高了性能。
func Mmap(fd int, offset int64, length int) ([]byte, error) {
	return syscall.Mmap(fd, offset, length, syscall.PROT_READ, syscall.MAP_SHARED)
}
