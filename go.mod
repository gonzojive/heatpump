module github.com/gonzojive/heatpump

go 1.14

//replace github.com/goburrow/serial => /home/pi/code/serial

// replace github.com/goburrow/modbus => /home/pi/code/modbus

require (
	github.com/dgraph-io/badger/v3 v3.2011.1
	github.com/fullstorydev/grpcui v1.1.0 // indirect
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.0
	github.com/google/go-cmp v0.5.5
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/inhies/go-bytesize v0.0.0-20201103132853-d0aed0d254f8
	github.com/martinlindhe/unit v0.0.0-20201217003049-aef7d8d7910f
	github.com/yuin/goldmark v1.3.0
	go.uber.org/multierr v1.6.0
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20201223074533-0d417f636930 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.26.0
)
