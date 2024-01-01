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
	"fmt"
	"os"
	"testing"

	"github.com/RoaringBitmap/roaring"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
)

func TestInvertedDBRead(t *testing.T) {
	query := "电影"
	termName := config.Conf.SeConfig.StoragePath + "0.term"
	inverted := NewInvertedDB(termName)
	v, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v", string(v))
	err = inverted.StoragePostings(query, []byte("100"))
	if err != nil {
		fmt.Println(err)
	}
	v2, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("v2", string(v2))
	err = inverted.PutInverted([]byte(query), []byte("11111"))
	if err != nil {
		fmt.Println(err)
	}
	v3, err := inverted.GetInverted([]byte(query))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(v3))
}

func TestStoreInvertedInfo(t *testing.T) {
	query := "蜘蛛侠"
	output := roaring.New()
	output.AddInt(1)
	output.AddInt(2)
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/mr-tmp-%d.%s",
		dir, 2, consts.InvertedBucket)
	inverted := NewInvertedDB(outName)
	oByte, _ := output.MarshalBinary()
	err := inverted.StoragePostings(query, oByte)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetInvertedInfo(t *testing.T) {
	query := "小岛"
	for i := 0; i < 5; i++ {
		outName := fmt.Sprintf("/Users/mac/GolandProjects/Go-SearchEngine/app/index_platform/woker/mr-tmp-%d.inverted", i)
		inverted := NewInvertedDB(outName)
		oByte, err := inverted.GetInverted([]byte(query))
		if err != nil {
			fmt.Println(err)
		}
		output := roaring.New()
		err = output.UnmarshalBinary(oByte)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
	}
}
