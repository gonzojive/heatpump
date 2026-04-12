package main

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/db"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/gonzojive/heatpump/tempsensor"
	"github.com/spf13/cobra"
)

var (
	histRoot           string
	histRs484TTYModbus string
	histDbDir          string
	histRestorePath    string
	histGrpcPort       int
)

const (
	reportInterval   = time.Minute
	snapshotInterval = time.Second * 5
)

var historianCmd = &cobra.Command{
	Use:   "start-historian",
	Short: "Start the historical collector service for the CX34",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		glog.Infof("starting start-historian")
		dbInst, err := db.Open(histDbDir)
		if err != nil {
			return err
		}
		if histRestorePath != "" {
			if err := dbInst.RestoreBackup(histRestorePath); err != nil {
				return err
			}
		}
		cxClient, err := cx34.Connect(&cx34.Params{TTYDevice: histRs484TTYModbus, Mode: cx34.Modbus})
		if err != nil {
			return err
		}

		config := &tempsensor.Config{Root: histRoot}
		eg, ctx := errgroup.WithContext(ctx)
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
					if err := dbInst.WriteSnapshot(state.Proto()); err != nil {
						glog.Errorf("error writing CX34 snapshot to database: %v", err)
						continue
					}
				}
			}
		})
		eg.Go(func() error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", histGrpcPort))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}
			s := grpc.NewServer()
			chiltrix.RegisterHistorianServer(s, dbInst.HistorianService())
			chiltrix.RegisterReadWriteServiceServer(s, cxClient.ReadWriteServiceServer())
			if err := s.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}
			return nil
		})
		return eg.Wait()
	},
}

func init() {
	historianCmd.Flags().StringVar(&histRoot, "root", "/", "Raspberry Pi filesystem root. Set to non-default value for testing with sshd on another computer.")
	historianCmd.Flags().StringVar(&histRs484TTYModbus, "modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	historianCmd.Flags().StringVar(&histDbDir, "db-dir", "/home/pi/db/cx34db", "Path to a directory where badger database should be stored.")
	historianCmd.Flags().StringVar(&histRestorePath, "restore-from-backup", "", "Path to a file with backup records when restoring from backup.")
	historianCmd.Flags().IntVar(&histGrpcPort, "grpc-port", 8082, "Port used to serve historical database values over GRPC.")
	rootCmd.AddCommand(historianCmd)
}
