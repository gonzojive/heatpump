package main

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/spf13/cobra"
)

var (
	rwRs484TTYModbus     string
	rwGrpcPort           int
	rwPrintStateInterval time.Duration
)

var readwriteCmd = &cobra.Command{
	Use:   "start-readwrite-service",
	Short: "Start a gRPC service bridging to the CX34 serial controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		glog.Infof("starting start-readwrite-service")
		cxClient, err := cx34.Connect(&cx34.Params{TTYDevice: rwRs484TTYModbus, Mode: cx34.Modbus})
		if err != nil {
			return err
		}
		eg, ctx := errgroup.WithContext(ctx)

		if rwPrintStateInterval.Seconds() != 0 {
			eg.Go(func() error {
				ticker := time.NewTicker(rwPrintStateInterval)
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
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", rwGrpcPort))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}
			s := grpc.NewServer()
			chiltrix.RegisterReadWriteServiceServer(s, cxClient.ReadWriteServiceServer())
			reflection.Register(s)
			glog.Infof("readwrite service server listening on :%d", rwGrpcPort)
			if err := s.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}
			return nil
		})
		return eg.Wait()
	},
}

func init() {
	readwriteCmd.Flags().StringVar(&rwRs484TTYModbus, "modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	readwriteCmd.Flags().IntVar(&rwGrpcPort, "grpc-port", 8084, "Port used to serve historical database values over GRPC.")
	readwriteCmd.Flags().DurationVar(&rwPrintStateInterval, "print-state-interval", 0, "If non-zero, the interval at which to print the state of the CX34 to the log.")
	rootCmd.AddCommand(readwriteCmd)
}
