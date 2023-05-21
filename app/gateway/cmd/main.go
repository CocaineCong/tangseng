package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"

	"github.com/CocaineCong/Go-SearchEngine/app/gateway/routes"
	"github.com/CocaineCong/Go-SearchEngine/config"
	"github.com/CocaineCong/Go-SearchEngine/pkg/discovery"
	"github.com/CocaineCong/Go-SearchEngine/pkg/util/shutdown"
)

func main() {
	config.InitConfig()
	// etcd注册
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewResolver(etcdAddress, logrus.New())
	defer etcdRegister.Close()
	resolver.Register(etcdRegister)
	go startListen() // 转载路由
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignals
		fmt.Println("exit! ", s)
	}
	fmt.Println("gateway listen on :3000")
}

func startListen() {
	ginRouter := routes.NewRouter()
	server := &http.Server{
		Addr:           config.Conf.Server.Port,
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("绑定HTTP到 %s 失败！可能是端口已经被占用，或用户权限不足 \n", config.Conf.Server.Port)
		fmt.Println(err)
	}
	go func() {
		// 优雅关闭
		shutdown.GracefullyShutdown(server)
	}()
}
