MODULE := $(shell awk '/^module / {print $$2}' go.mod)
PB_FILES:=$(shell find . -type f -name '*.proto' -print)
GOPB_FILES:=$(PB_FILES:%.proto=%.pb.go)

all: $(GOPB_FILES) bin/server bin/client

$(GOPB_FILES): $(PB_FILES)
	@protoc \
		--go_out=module=$(MODULE):. \
		--go-grpc_out=module=$(MODULE):. \
		--doc_out=. --doc_opt=html,proto_doc.html \
		$(@:%.pb.go=%.proto)

bin/server: $(GOPB_FILES) $(wildcard server/*.go)
	go build -o $@ ./server

bin/client: $(GOPB_FILES) $(wildcard client/*.go)
	go build -o $@ ./client

.PHONY: clean
clean:
	rm -rf *.html proto bin
