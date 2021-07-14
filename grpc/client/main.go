package main

import (
	"github.com/care0717/deepthought-api/grpc/client/command"
	"log"
)

func main() {
	if err := command.Run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
