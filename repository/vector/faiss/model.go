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

package faiss

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/config"
)

type FaissModel struct {
	name   string
	client *VectorClient
}

func NewFaissModel(name string) *FaissModel {
	return &FaissModel{
		name: name,
	}
}

func (m *FaissModel) Init(ctx context.Context) (err error) {
	vConfig := config.Conf.Vector
	client, err := NewVectorClient(ctx, vConfig.ServerAddress, time.Millisecond*time.Duration(vConfig.Timeout))
	if err != nil {
		return errors.Wrap(err, "failed to create new vector client")
	}
	m.client = client

	return
}

func (m *FaissModel) Run(data interface{}) (resp interface{}, err error) {
	resp, err = m.client.Search(data)
	if err != nil {
		err = errors.WithMessage(err, "search error")
	}
	return
}
