# Project protocol buffers

To regenerate the .pb.go file, run this from the root of the project:

```shell
protoc -I=proto/chiltrix \
  --go_opt=paths=source_relative \
  --go_out=proto/chiltrix \
  --go-grpc_opt=paths=source_relative \
  --go-grpc_out=proto/chiltrix \
  proto/chiltrix/chiltrix.proto
```


## Prerequisites

```
sudo apt install protobuf-compiler
```

```
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
