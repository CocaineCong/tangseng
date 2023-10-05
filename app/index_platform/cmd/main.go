package main

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/job"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/kfk_register"
	"github.com/CocaineCong/tangseng/app/index_platform/service"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const (
	IndexPlatformServerName = "index_platform"
)

func main() {
	ctx := context.Background()
	// 加载配置
	loading.Loading()
	analyzer.InitSeg()
	kfk_register.RegisterJob(ctx)
	job.RegisterJob(ctx)

	// 注册服务
	_ = registerIndexPlatform()
}

// registerIndexPlatform 注册索引平台服务
func registerIndexPlatform() (err error) {
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	defer etcdRegister.Stop()

	grpcAddress := config.Conf.Services[IndexPlatformServerName].Addr[0]
	node := discovery.Server{
		Name: config.Conf.Domain[IndexPlatformServerName].Name,
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
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}

	return
}
