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

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/service/recall"
	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/repository/mysql/db"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	analyzer.InitSeg()
	log.InitLog()
	db.InitDB()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestRecall(t *testing.T) {
	q := "小岛"
	ctx := context.Background()
	searchItem, err := recall.SearchRecall(ctx, q)
	if err != nil {
		fmt.Println(err)
	}
	for i := range searchItem {
		fmt.Println(searchItem[i].Score, searchItem[i].DocId, searchItem[i].Content)
	}
}
