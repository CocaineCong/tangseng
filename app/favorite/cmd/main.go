package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/favorite/internal/service"
	"github.com/CocaineCong/tangseng/app/gateway/rpc"
	"github.com/CocaineCong/tangseng/config"
	favoritePb "github.com/CocaineCong/tangseng/idl/pb/favorite"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

const ServerName = "favorite"

func main() {
	loading.Loading()
	rpc.Init()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logs.LogrusObj)
	grpcAddress := config.Conf.Services[ServerName].Addr[0]
	defer etcdRegister.Stop()
	userNode := discovery.Server{
		Name: config.Conf.Domain[ServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	favoritePb.RegisterFavoritesServiceServer(server, service.GetFavoriteSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
