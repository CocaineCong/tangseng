package prometheus

import (
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

var (
	UnaryServerInterceptor  = grpcPrometheus.UnaryServerInterceptor
	StreamServerInterceptor = grpcPrometheus.StreamServerInterceptor

	EnableHandlingTimeHistogram = grpcPrometheus.EnableHandlingTimeHistogram
)

// RegisterServer  for prometheus
func RegisterServer(server *grpc.Server, domain string, job string) {
	grpcPrometheus.Register(server)
	EtcdRegister(domain, job)
	go RpcHandler(domain)
}
