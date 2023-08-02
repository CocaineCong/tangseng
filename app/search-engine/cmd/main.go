package main

import (
	"github.com/CocaineCong/tangseng/app/search-engine/internal/index"
	"github.com/CocaineCong/tangseng/app/search-engine/internal/query"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/loading"
)

func main() {
	config.InitConfig()
	loading.Loading()
	query.InitSeg()
	index.RunningIndex()
	// query := "孙悟空"
	// index.SearchRecall(query)

	// // etcd 地址
	// etcdAddress := []string{viper.GetString("etcd.address")}
	// // 服务注册
	// etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	// grpcAddress := viper.GetString("server.grpcAddress")
	// defer etcdRegister.Stop()
	// userNode := discovery.Server{
	// 	Name: viper.GetString("server.domain"),
	// 	Addr: grpcAddress,
	// }
	// server := grpc.NewServer()
	// defer server.Stop()
	// // 绑定service
	// pb.RegisterSearchEngineServiceServer(server, handler.GetSearchEngineSrv())
	// lis, err := net.Listen("tcp", grpcAddress)
	// if err != nil {
	// 	panic(err)
	// }
	// if _, err := etcdRegister.Register(userNode, 10); err != nil {
	// 	panic(fmt.Sprintf("start server failed, err: %v", err))
	// }
	// logrus.Info("server started listen on ", grpcAddress)
	// if err := server.Serve(lis); err != nil {
	// 	panic(err)
	// }
}
