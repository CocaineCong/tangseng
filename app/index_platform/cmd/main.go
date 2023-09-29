package main

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/index_platform/analyzer"
	"github.com/CocaineCong/tangseng/app/index_platform/cmd/kfk_register"
	"github.com/CocaineCong/tangseng/app/index_platform/repository/storage"
	"github.com/CocaineCong/tangseng/app/index_platform/rpc"
	"github.com/CocaineCong/tangseng/app/index_platform/service"
	"github.com/CocaineCong/tangseng/app/index_platform/trie"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const (
	IndexPlatformServerName = "index_platform"
	MapreduceServerName     = "mapreduce"
)

func main() {
	ctx := context.Background()
	loading.Loading()
	analyzer.InitSeg()
	rpc.Init()
	trie.InitTrieTree()
	storage.InitTrieDBs()
	kfk_register.RegisterJob(ctx)

	go func() {
		err := registerMapreduce()
		if err != nil {
			logs.LogrusObj.Errorln("registerMapreduce", err)
		}
	}()
	// _ = registerMapreduce()
	_ = registerIndexPlatform()
}

// registerMapreduce 注册mapreduce服务
func registerMapreduce() (err error) {
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

	mapreduce.RegisterMapReduceServiceServer(server, service.GetMapReduceSrv())
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
