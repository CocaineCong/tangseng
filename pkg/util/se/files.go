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

package se

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Walk(dirPath string) []string {
	var fileList []string
	_ = filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	// fmt.Println(len(fileList))

	return fileList
}

func GetMd5(filePath string) string {
	file, err := os.Open(filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)
	if err != nil {
		panic(err)
	}

	md5h := md5.New()
	_, _ = io.Copy(md5h, file)
	return hex.EncodeToString(md5h.Sum(nil))
}

func DirCHeckAndMk(dir string) {
	_, err := os.Stat(dir)
	if err != nil { // 文件不存在
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("permission denied![%v]\n", err)
			panic(err)
		}
	}
}

// ExistFile 判断所给的路径文件/文件夹是否存在
func ExistFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
