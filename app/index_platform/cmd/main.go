package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/rpc"
	"github.com/CocaineCong/tangseng/app/index_platform/service"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const ServerName = "index_platform"

func main() {
	// loading.Loading()
	config.InitConfig()
	logs.InitLog()
	analyzer.InitSeg()
	rpc.Init()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	grpcAddress := config.Conf.Services[ServerName].Addr[0]
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: config.Conf.Domain[ServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	mapreduce.RegisterMapReduceServiceServer(server, service.GetMapReduceSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
