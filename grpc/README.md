# grpc-tutorial for deepthrought
https://github.com/ymmt2005/grpc-tutorial の実装

## require
- protoc
- protoc-gen-go
- protoc-gen-doc
- protoc-gen-grpc

## usage
### build
```shell
make
```

### run server
```shell
./bin/server
```

### run client
```shell
./bin/clinet -h
./bin/clinet boot
./bin/clinet infer Life
```

## read proto docs
```shell
make deepthought.html
open deepthought.html
```
