package main

import (
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	portNumber = 13333
	promPort   = 18888
)

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	serv := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kep),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)
	deepthought.RegisterComputeServer(serv, &Server{})
	grpc_prometheus.Register(serv)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Info(fmt.Sprintf("prometheus metrics bind port: %d", promPort))
		logger.Fatal(fmt.Sprintf("listen failed: %v", http.ListenAndServe(fmt.Sprintf(":%d", promPort), mux)))
	}()
	// 待ち受けソケットを作成
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	logger.Info(fmt.Sprintf("listen port :%d", portNumber))

	serv.Serve(l)
}
