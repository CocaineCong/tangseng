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

package milvus

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"

	"github.com/CocaineCong/tangseng/config"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type MilvusModel struct {
	ctx    context.Context
	name   string
	client client.Client
}

func NewMilvusModel(ctx context.Context, name string) *MilvusModel {
	return &MilvusModel{ctx: ctx, name: name}
}

func (m *MilvusModel) Init() (err error) {
	mConfig := config.Conf.Milvus
	ctx, cancel := context.WithTimeout(m.ctx, time.Millisecond*time.Duration(mConfig.Timeout))
	defer cancel()
	milvusClient, err := client.NewGrpcClient(ctx, fmt.Sprintf("%s:%s", mConfig.Host, mConfig.Port))
	if err != nil {
		logs.LogrusObj.Errorln(err)
		return
	}
	m.client = milvusClient

	return
}

func (m *MilvusModel) Search(req interface{}) (resp interface{}, err error) {
	request, ok := req.(*MilvusRequest)
	if !ok {
		return
	}

	return m.client.Search(
		m.ctx,
		request.CollectionName,
		request.Partitions,
		request.Expr,
		request.OutputFields,
		request.Vectors,
		request.VectorField,
		request.MetricType,
		request.TopK,
		request.SearchParams,
		nil,
	)
}
