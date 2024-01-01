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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/CocaineCong/tangseng/idl/pb/vector_retrieval"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

type VectorClient struct {
	ctx           context.Context
	ServerAddress string
	Timeout       time.Duration
	VectorClient  pb.VectorRetrievalClient
}

func NewVectorClient(ctx context.Context, address string, timeout time.Duration) (client *VectorClient, err error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	client = &VectorClient{
		ctx:           ctx,
		ServerAddress: address,
		Timeout:       timeout,
		VectorClient:  pb.NewVectorRetrievalClient(conn),
	}

	return
}

func (c *VectorClient) Search(req interface{}) (resp *pb.VectorResp, err error) {
	request, ok := req.(*pb.VectorReq)
	if !ok {
		return
	}
	ctx, cancl := context.WithTimeout(c.ctx, c.Timeout)
	defer cancl()
	resp, err = c.VectorClient.Search(ctx, request)
	if err != nil {
		logs.LogrusObj.Errorln("VectorClient-Search ", err)
		return
	}

	return
}
