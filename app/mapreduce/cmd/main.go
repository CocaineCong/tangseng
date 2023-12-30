package main

import (
	"github.com/pkg/errors"
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
		logs.LogrusObj.Errorf("start service failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		logs.LogrusObj.Panicf("stack trace: \n%+v\n", err)
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
