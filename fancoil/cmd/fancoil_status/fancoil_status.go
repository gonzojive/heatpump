// Program fancoil_status prints out the status of any connected fan coil.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/fancoil"
	fcpb "github.com/gonzojive/heatpump/proto/fancoil"
	"google.golang.org/grpc"
)

var (
	rs484TTYModbus = flag.String("modbus-device", "/dev/serial/by-id/usb-1a86_USB_Serial-if00-port0", "Path to USB-to-RS485 device connected to modbus.")
	startServer    = flag.Bool("start-server", false, "Start a gRPC server.")
	grpcPort       = flag.Int("grpc-port", 8082, "Port used to serve historical database values over GRPC.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}
func run(ctx context.Context) error {
	client, err := fancoil.Connect(ctx, &fancoil.Params{
		TTYDevice: *rs484TTYModbus,
	})
	if err != nil {
		return fmt.Errorf("error connecting: %w", err)
	}

	if *startServer {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s := grpc.NewServer()
		fcpb.RegisterFanCoilServiceServer(s, fancoil.GRPCServiceFromClient(client))
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	}
	return nil
}
