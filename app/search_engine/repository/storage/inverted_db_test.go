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

package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/CocaineCong/tangseng/config"
)

func TestGetInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	postingsName := config.Conf.SeConfig.StoragePath + "0.inverted"
	inverted := NewInvertedDB(termName, postingsName)
	p, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
}

func TestInitInvertedDB(t *testing.T) {
	ctx := context.Background()
	InitInvertedDB(ctx)
	for _, v := range GlobalInvertedDB {
		fmt.Println(v)
	}
}
