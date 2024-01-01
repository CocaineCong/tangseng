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

package kfk

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../config/config.yaml"}
	config.InitConfigForTest(&re)
	InitKafka()
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

type TestKafkaData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TestKafkaProducer(t *testing.T) {
	data := &TestKafkaData{
		Key:   "怎么说",
		Value: "啊哈哈哈哈",
	}
	d, _ := json.Marshal(data)
	err := KafkaProducer(consts.KafkaCSVLoaderTopic, d)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Produce Message Finish")
}
