# Project protocol buffers

To regenerate the .pb.go file, run this from the root of the project:

```shell
protoc -I=proto/chiltrix --go_opt=paths=source_relative --go_out=proto/chiltrix proto/chiltrix/chiltrix.proto
```
