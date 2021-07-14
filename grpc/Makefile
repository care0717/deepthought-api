MODULE := $(shell awk '/^module / {print $$2}' go.mod)

all: deepthought.html deepthought.pb.go deepthought_grpc.pb.go bin/server bin/client

deepthought.html: deepthought.proto
	protoc --doc_out=. --doc_opt=html,$@ $<

deepthought.pb.go: deepthought.proto
	protoc --go_out=module=$(MODULE):. $<

deepthought_grpc.pb.go: deepthought.proto
	protoc --go-grpc_out=module=$(MODULE):. $<

bin/server: deepthought.pb.go deepthought_grpc.pb.go $(wildcard server/*.go)
	go build -o $@ ./server

bin/client: deepthought.pb.go deepthought_grpc.pb.go $(wildcard client/*.go)
	go build -o $@ ./client

.PHONY: clean
clean:
	rm -rf deepthought.index proto bin
