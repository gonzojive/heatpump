// Program cx34collector collects information from the Chiltrix CX34 heat pump
// and logs it to a database.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/proto/chiltrix"
)

var (
	root               = flag.String("root", "/", "Raspberry Pi filesystem root. Set to non-default value for testing with sshd on another computer.")
	rs484TTYModbus     = flag.String("modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	grpcPort           = flag.Int("grpc-port", 8083, "Port used to serve historical database values over GRPC.")
	versionFlag        = flag.Bool("version", false, "Return the version of the program.")
	printStateInterval = flag.Duration("print-state-interval", time.Hour*0, "If non-zero, the interval at which to print the state of the CX34 to the log.")

	markdown = goldmark.New(goldmark.WithExtensions(extension.NewTable()))
)

const (
	reportInterval   = time.Minute
	snapshotInterval = time.Second * 5

	version = "v0.0.2"
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Printf("%s\n", version)
		return
	}
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	glog.Infof("starting cx34control version %s", version)
	cxClient, err := cx34.Connect(&cx34.Params{TTYDevice: *rs484TTYModbus, Mode: cx34.Modbus})
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(context.Background())

	if printStateInterval.Seconds() != 0 {
		eg.Go(func() error {
			ticker := time.NewTicker(*printStateInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					state, err := cxClient.ReadState()
					if err != nil {
						glog.Errorf("error getting CX34 state: %v", err)
						continue
					}
					glog.Infof("State of CX34 heat pump: %s", state.Report(false, nil))
				}
			}
		})
	}

	eg.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s := grpc.NewServer()
		chiltrix.RegisterReadWriteServiceServer(s, cxClient.ReadWriteServiceServer())
		reflection.Register(s)
		glog.Infof("cx34control server listening on :%d", *grpcPort)
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	})
	return eg.Wait()
}
