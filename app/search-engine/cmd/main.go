package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/CocaineCong/Go-SearchEngine/app/search-engine/internal/handler"
	"github.com/CocaineCong/Go-SearchEngine/config"
	pb "github.com/CocaineCong/Go-SearchEngine/idl/pb/search_engine"
	"github.com/CocaineCong/Go-SearchEngine/pkg/discovery"
)

func main() {
	config.InitConfig()
	// etcd 地址
	etcdAddress := []string{viper.GetString("etcd.address")}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := viper.GetString("server.grpcAddress")
	defer etcdRegister.Stop()
	userNode := discovery.Server{
		Name: viper.GetString("server.domain"),
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	pb.RegisterSearchEngineServiceServer(server, handler.GetSearchEngineSrv())
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
