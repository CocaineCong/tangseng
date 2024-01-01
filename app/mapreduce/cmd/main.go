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

package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/mapreduce/master"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const (
	MapreduceServerName = "mapreduce"
)

func main() {
	loading.Loading()
	analyzer.InitSeg()

	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	defer etcdRegister.Stop()

	grpcAddress := config.Conf.Services[MapreduceServerName].Addr[0]
	node := discovery.Server{
		Name: config.Conf.Domain[MapreduceServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	mapreduce.RegisterMapReduceServiceServer(server, master.GetMapReduceSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
