module github.com/gonzojive/heatpump

go 1.22.10

// replace github.com/rmrobinson/google-smart-home-action-go => github.com/gonzojive/google-smart-home-action-go v0.0.1

//replace github.com/goburrow/serial => /home/pi/code/serial

// replace github.com/goburrow/modbus => /home/pi/code/modbus

require (
	cloud.google.com/go/firestore v1.17.0
	cloud.google.com/go/pubsub v1.45.3
	cloud.google.com/go/secretmanager v1.14.2
	github.com/adrg/xdg v0.5.3
	github.com/bazelbuild/rules_go v0.51.0
	github.com/dgraph-io/badger/v3 v3.2103.5
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/go-oauth2/oauth2/v4 v4.5.2
	github.com/goburrow/modbus v0.1.0
	github.com/goburrow/serial v0.1.0
	github.com/golang/glog v1.2.3
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.6.0
	github.com/google/go-containerregistry v0.20.2
	github.com/howeyc/crc16 v0.0.0-20171223171357-2b2a61e366a6
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf
	github.com/johnsiilver/golib v1.2.2
	github.com/martinlindhe/unit v0.0.0-20230420213220-4adfd7d0a0d6
	github.com/mtraver/iotcore v0.0.0-20230423225757-7fc79eb2c3c4
	github.com/rmrobinson/google-smart-home-action-go v0.0.0-20240904013938-6a5c976efa23
	github.com/samber/lo v1.47.0
	github.com/stretchr/testify v1.10.0
	github.com/teambition/rrule-go v1.8.2
	github.com/yuin/goldmark v1.7.8
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.31.0
	golang.org/x/net v0.33.0
	golang.org/x/oauth2 v0.24.0
	golang.org/x/sync v0.10.0
	google.golang.org/api v0.214.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.1
)

require (
	cloud.google.com/go v0.117.0 // indirect
	cloud.google.com/go/auth v0.13.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.3.0 // indirect
	cloud.google.com/go/longrunning v0.6.3 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.16.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/docker/cli v27.4.1+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.8.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v24.12.23+incompatible // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tidwall/buntdb v1.3.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/vbatts/tar-split v0.11.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.58.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/dig v1.18.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20241223144023-3abc09e42ca8 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241223144023-3abc09e42ca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241223144023-3abc09e42ca8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
