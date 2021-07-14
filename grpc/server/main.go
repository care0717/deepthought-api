package main

import (
	"fmt"
	"github.com/care0717/deepthought-api/grpc/proto/deepthought"
	"google.golang.org/grpc"
	"net"
	"os"
)

const portNumber = 13333

func main() {
	serv := grpc.NewServer()

	deepthought.RegisterComputeServer(serv, &Server{})

	// 待ち受けソケットを作成
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	serv.Serve(l)
}
