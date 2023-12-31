package main

import (
	"context"
	"net"

	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
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
)

func main() {
	ctx := context.Background()
	loading.Loading()
	// bi_dao.InitDB() // TODO starrocks完善才开启
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
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	pb.RegisterSearchEngineServiceServer(server, service.GetSearchEngineSrv())
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
