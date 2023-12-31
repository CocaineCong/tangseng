package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/idl/pb/search_vector"
	"github.com/CocaineCong/tangseng/pkg/discovery"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	SearchVectorClient search_vector.SearchVectorServiceClient
)

// Init 初始化所有的rpc请求
func Init() {
	fmt.Println("loading rpc init")
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	initClient(config.Conf.Domain[consts.SearchVectorName].Name, &SearchVectorClient)
	fmt.Println("loading rpc init finished")
}

// initClient 初始化所有的rpc客户端
func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		logs.LogrusObj.Errorf("start service failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		panic(err)
	}

	switch c := client.(type) {
	case *search_vector.SearchVectorServiceClient:
		*c = search_vector.NewSearchVectorServiceClient(conn)
	default:
		panic("unsupported worker type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	// Load balance
	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	err = errors.Wrap(err, "failed to connect grpc")
	return
}
