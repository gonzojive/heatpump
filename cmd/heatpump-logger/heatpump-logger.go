// Program cx34runs a dashboard that displays information about a CX34 heat pump.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/protobuf/proto"

	"github.com/gonzojive/heatpump/cmd/heatpump-logger/loglib"
	"github.com/gonzojive/heatpump/grpcspec"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"go.uber.org/fx"

	fancoilpb "github.com/gonzojive/heatpump/proto/fancoil"
)

var (
	httpPort           = flag.Int("status-http-port", 8081, "HTTP server port for reporting the logger's status.")
	dataDir            = flag.String("data-dir", "", "Directory where .")
	newLogFileInterval = flag.Duration("log-file-interval", time.Hour, "Maximum amount of time between initial log file creation and its finalization.")

	readWriteServiceParams = grpcspec.ClientParamsFlag{
		Value: grpcspec.NewClientParamsBuilder().
			SetAddr("localhost:8084").
			SetInsecure(true).
			Build(),
	}

	fancoilServiceParams = grpcspec.ClientParamsFlag{
		Value: grpcspec.NewClientParamsBuilder().
			SetAddr("localhost:8083").
			SetInsecure(true).
			Build(),
	}
)

func init() {
	flag.Var(&readWriteServiceParams, "read-write-grpc-service", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true) for chiltrix.ReadWriteService")
	flag.Var(&fancoilServiceParams, "fan-coil-grpc-service", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true) for FanCoilService")
}

func main() {
	flag.Parse()
	fx.New(
		fx.Provide(newHTTPServer),
		fx.Provide(newReadWriteServiceClient),
		fx.Provide(newFancoilServiceClient),
		fx.Provide(newLoggingService),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func newReadWriteServiceClient(lc fx.Lifecycle) (chiltrix.ReadWriteServiceClient, error) {
	client, conn, err := grpcspec.NewClient(readWriteServiceParams.Value, chiltrix.NewReadWriteServiceClient)
	if err != nil {
		return nil, fmt.Errorf("error creating ReadWriteServiceClient: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return client, nil
}

func newFancoilServiceClient(lc fx.Lifecycle) (fancoilpb.FanCoilServiceClient, error) {
	client, conn, err := grpcspec.NewClient(fancoilServiceParams.Value, fancoilpb.NewFanCoilServiceClient)
	if err != nil {
		return nil, fmt.Errorf("error creating FanCoilServiceClient: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return client, nil
}

type LoggingService struct {
	logWriter *loglib.MultiFileTFRecordWriter
}

func periodicallyLogHeatpumpState(
	ctx context.Context,
	interval time.Duration,
	done <-chan struct{},
	logger *LoggingService,
	heatpumpClient chiltrix.ReadWriteServiceClient,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			state, err := heatpumpClient.GetState(ctx, &chiltrix.GetStateRequest{})
			if err != nil {
				glog.Errorf("error getting heatpump state: %v", err)
			}
			registerSnapshotBytes, err := proto.Marshal(state.GetRegisterValues())
			if err != nil {
				glog.Errorf("error encoding RegisterValues proto: %v", err)
			}
			if err := logger.logWriter.Write(registerSnapshotBytes); err != nil {
				glog.Errorf("error writing RegisterValues proto: %v", err)
			}
		case <-done:
			return
		}
	}
}

func newLoggingService(lc fx.Lifecycle, heatpumpClient chiltrix.ReadWriteServiceClient) *LoggingService {
	service := &LoggingService{}
	doneCh := make(chan struct{})
	wg := &sync.WaitGroup{}
	var loggerCtx context.Context
	var cancelLoggerCtx context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if *dataDir == "" {
				return fmt.Errorf("invalid --data-dir argument (blank)")
			}
			loggerCtx, cancelLoggerCtx = context.WithCancel(context.Background())
			service.logWriter = loglib.NewPeriodicMultiFileTFRecordWriter(
				loggerCtx,
				time.Now,
				filepath.Join(*dataDir, "cx34-RegisterSnapshot"),
				".tfrecord",
				*newLogFileInterval,
			)

			go func() {
				wg.Add(1)
				defer wg.Done()
				periodicallyLogHeatpumpState(loggerCtx, time.Second*10, doneCh, service, heatpumpClient)
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancelLoggerCtx()
			wg.Wait()
			return service.logWriter.Close()
		},
	})
	return service
}

func newHTTPServer(lc fx.Lifecycle, loggingService *LoggingService) *http.Server {
	mux := http.NewServeMux()
	// Add your handler for /index.md here
	mux.HandleFunc("/index.md", func(w http.ResponseWriter, r *http.Request) {
		markdown := ""
		w.Header().Set("Content-Type", "text/markdown")
		fmt.Fprint(w, markdown)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
