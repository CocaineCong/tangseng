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

package starrock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

// mysql -h 127.0.0.1 -P9030 -uroot
func TestDirectUpload_StreamUpload(t *testing.T) {
	ctx := context.Background()
	du := NewDirectUpload(ctx, &types.Task{
		Columns:    []string{"doc_id", "url", "title", "desc", "score"},
		BiTable:    "test_upload",
		SourceType: 0,
	})
	fmt.Println(config.Conf.StarRocks)
	du.Push(&types.Data2Starrocks{
		DocId: 1,
		Url:   "https://localhost:8083",
		Title: "这是一个测试文件",
		Desc:  "进行测试作用",
		Score: 1220.120,
	})
	time.Sleep(10 * time.Second)
}
