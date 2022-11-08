// Program authservice starts the AuthService. It is intended to be hosted in
// Google Cloud.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/util/grpcserverutil"

	pb "github.com/gonzojive/heatpump/proto/controller"
)

var (
	grpcPort = flag.Int("grpc-port", 8092, "Port used to serve gRPC traffic.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	glog.Infof("starting gRPC service at :%d", *grpcPort)

	var grpcOpts []grpc.ServerOption
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s := grpc.NewServer(grpcOpts...)
		authService, err := acls.NewAuthService(ctx, cloudconfig.DefaultParams())
		if err != nil {
			return err
		}
		pb.RegisterAuthServiceServer(s, authService)
		reflection.Register(s)

		if err := grpcserverutil.ServeUntilCancelled(ctx, s, lis); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	})
	return eg.Wait()
}
