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

```shell
protoc -I=proto/controller -I=.   --go_opt=paths=source_relative   --go_out=proto/controller --go-grpc_opt=paths=source_relative   --go-grpc_out=proto/controller   proto/controller/controller.proto
```

```shell
protoc -I=proto/command_queue   --go_opt=paths=source_relative   --go_out=proto/command_queue   --go-grpc_opt=paths=source_relative   --go-grpc_out=proto/command_queue   proto/command_queue/command_queue.proto
```


## Prerequisites

```
sudo apt install protobuf-compiler
```

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
