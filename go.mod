module github.com/gonzojive/heatpump

go 1.14

//replace github.com/goburrow/serial => /home/pi/code/serial

//replace github.com/goburrow/modbus => /home/pi/code/modbus

require (
	github.com/dgraph-io/badger/v2 v2.0.3
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.0
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/martinlindhe/unit v0.0.0-20201217003049-aef7d8d7910f
	github.com/yuin/goldmark v1.3.0
	go.uber.org/multierr v1.6.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/grpc v1.33.0-dev
	google.golang.org/protobuf v1.25.0
)
