package prometheus

import (
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

var (
	// UnaryServerInterceptor intercept all unary rpc requests and set metrics
	UnaryServerInterceptor = grpcPrometheus.UnaryServerInterceptor
	// StreamServerInterceptor intercept all stream rpc requests and set metrics
	StreamServerInterceptor = grpcPrometheus.StreamServerInterceptor

	// EnableHandlingTimeHistogram enable the function of handling time histogram
	EnableHandlingTimeHistogram = grpcPrometheus.EnableHandlingTimeHistogram
)

// RegisterServer for prometheus
func RegisterServer(server *grpc.Server, domain string, job string) {
	grpcPrometheus.Register(server)
	EtcdRegister(domain, job)
	go RpcHandler(domain)
}
