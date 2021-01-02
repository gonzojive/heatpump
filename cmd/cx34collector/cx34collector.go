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

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/db"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/gonzojive/heatpump/tempsensor"
)

var (
	root           = flag.String("root", "/", "Raspberry Pi filesystem root. Set to non-default value for testing with sshd on another computer.")
	rs484TTYModbus = flag.String("modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	dbDir          = flag.String("db-dir", "/home/pi/db/cx34db", "Path to a directory where badger database should be stored.")
	grpcPort       = flag.Int("grpc-port", 8082, "Port used to serve historical database values over GRPC.")
	versionFlag    = flag.Bool("version", false, "Return the version of the program.")

	markdown = goldmark.New(goldmark.WithExtensions(extension.NewTable()))
)

const (
	reportInterval   = time.Minute
	snapshotInterval = time.Second * 5

	version = "v0.0.1a"
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
	db, err := db.Open(*dbDir)
	if err != nil {
		return err
	}
	cxClient, err := cx34.Connect(&cx34.Params{TTYDevice: *rs484TTYModbus, Mode: cx34.Modbus})
	if err != nil {
		return err
	}

	config := &tempsensor.Config{Root: *root}
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		ticker := time.NewTicker(reportInterval)

		defer ticker.Stop()
		for {
			r, err := tempsensor.DebugReport(config)
			if err != nil {
				glog.Errorf("got error: %v", err)
			} else {
				glog.Infof("%s\n", r)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
	})

	eg.Go(func() error {
		ticker := time.NewTicker(snapshotInterval)
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
				if err := db.WriteSnapshot(state.Proto()); err != nil {
					glog.Errorf("error writing CX34 snapshot to database: %v", err)
					continue
				}
			}
		}
	})
	eg.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s := grpc.NewServer()
		chiltrix.RegisterHistorianServer(s, db.HistorianService())
		chiltrix.RegisterReadWriteServiceServer(s, cxClient)
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	})
	return eg.Wait()
}
