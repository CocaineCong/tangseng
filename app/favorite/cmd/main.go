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
)

func main() {
	loading.Loading()
	rpc.Init()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["favorite"].Addr[0]
	defer etcdRegister.Stop()
	userNode := discovery.Server{
		Name: config.Conf.Domain["favorite"].Name,
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
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
