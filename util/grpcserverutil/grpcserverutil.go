// Package grpcserverutil contains utilities for running gRPC servers.
package grpcserverutil

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// ServeUntilCancelled calls server.Serve until the provided context is
// cancelled. When cancellation occurs, server.GracefulStop is called.
func ServeUntilCancelled(ctx context.Context, server *grpc.Server, lis net.Listener) error {
	serverStopped := make(chan struct{})
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			server.GracefulStop()
			return ctx.Err()
		case <-serverStopped:
			return nil
		}
	})
	eg.Go(func() error {
		defer close(serverStopped)
		return server.Serve(lis)
	})
	return eg.Wait()
}
