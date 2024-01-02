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
	"context"
	"net"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/job"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/kfk_register"
	"github.com/CocaineCong/tangseng/app/index_platform/service"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

func main() {
	ctx := context.Background()
	// 加载配置
	loading.Loading()
	analyzer.InitSeg()
	kfk_register.RegisterJob(ctx)
	job.RegisterJob(ctx)

	// 注册服务
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	defer etcdRegister.Stop()
	grpcAddress := config.Conf.Services[consts.IndexPlatformName].Addr[0]
	node := discovery.Server{
		Name: config.Conf.Domain[consts.IndexPlatformName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	index_platform.RegisterIndexPlatformServiceServer(server, service.GetIndexPlatformSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(node, 10); err != nil {
		logs.LogrusObj.Errorf("start service failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		logs.LogrusObj.Panicf("stack trace: \n%+v\n", err)
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
