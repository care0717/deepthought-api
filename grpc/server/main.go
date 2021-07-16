package main

import (
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"time"
)

const portNumber = 13333

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 10 * time.Second,
	}
	serv := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kep))

	deepthought.RegisterComputeServer(serv, &Server{})

	// 待ち受けソケットを作成
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	serv.Serve(l)
}
