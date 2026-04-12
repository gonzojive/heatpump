module github.com/gonzojive/heatpump

go 1.23.4

// replace github.com/rmrobinson/google-smart-home-action-go => github.com/gonzojive/google-smart-home-action-go v0.0.1

//replace github.com/goburrow/serial => /home/pi/code/serial

// replace github.com/goburrow/modbus => /home/pi/code/modbus

require (
	github.com/bazelbuild/rules_go v0.51.0
	github.com/dgraph-io/badger/v3 v3.2103.5
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/golang/glog v1.2.3
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.6.0
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf
	github.com/martinlindhe/unit v0.0.0-20230420213220-4adfd7d0a0d6
	github.com/ryszard/tfutils v0.0.0-20161028141955-98de232c7c68
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/teambition/rrule-go v1.8.2
	github.com/yuin/goldmark v1.7.8
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.1
)

require (
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v24.12.23+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241223144023-3abc09e42ca8 // indirect
)
