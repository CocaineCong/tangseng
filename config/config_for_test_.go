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

package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// 这个文件是为了方便写的test文件来读取config

type IReader interface {
	readConfig() ([]byte, error)
}

type ConfigReader struct {
	FileName string
}

// 'reader' implementing the Interface
// Function to read from actual file
func (r *ConfigReader) readConfig() ([]byte, error) {
	file, err := os.ReadFile(r.FileName)

	if err != nil {
		log.Fatal(err)
	}
	return file, err
}

func InitConfigForTest(reader IReader) {
	file, err := reader.readConfig()
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		panic(err)
	}
}
