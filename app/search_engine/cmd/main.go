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
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/rpc"
	"github.com/CocaineCong/tangseng/app/search_engine/service"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/pkg/prometheus"
	"github.com/CocaineCong/tangseng/pkg/tracing"
)

func main() {
	ctx := context.Background()
	loading.Loading()
	// bi_dao.InitDB() // TODO: starrocks完善才开启
	analyzer.InitSeg()
	storage.InitStorageDB(ctx)
	rpc.Init()

	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services[consts.SearchServiceName].Addr[0]
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: config.Conf.Domain[consts.SearchServiceName].Name,
		Addr: grpcAddress,
	}
	// 注册tracer
	provider := tracing.InitTracerProvider(config.Conf.Jaeger.Addr, consts.SearchServiceName)
	defer func() {
		if provider == nil {
			return
		}
		if err := provider(context.Background()); err != nil {
			logs.LogrusObj.Errorf("Failed to shutdown: %v", err)
		}
	}()
	handler := otelgrpc.NewServerHandler()
	server := grpc.NewServer(
		grpc.StatsHandler(handler),
		grpc.UnaryInterceptor(prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(prometheus.StreamServerInterceptor),
	)

	defer server.Stop()
	// 绑定service
	pb.RegisterSearchEngineServiceServer(server, service.GetSearchEngineSrv())
	prometheus.RegisterServer(server, config.Conf.Services[consts.SearchServiceName].Metrics[0], consts.SearchServiceName)
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		logs.LogrusObj.Errorf("start service failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		logs.LogrusObj.Panicf("stack trace: \n%+v\n", err)
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
