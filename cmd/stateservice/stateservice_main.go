// Program stateservice starts a controller.StateService intende to be hosted in
// Google Cloud.
//
// The service uses FireStore to store the device state and runs under Cloud
// Run (typically).
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/stateservice"
	pb "github.com/gonzojive/heatpump/proto/controller"
	"github.com/gonzojive/heatpump/util/grpcserverutil"
)

var (
	grpcPort             = flag.Int("grpc-port", 8089, "Port used to serve gRPC traffic.")
	insecure             = flag.Bool("insecure", true, "If true, don't load a server TLS certificate.")
	serverCertPath       = flag.String("server-cert", "", "Path to server cert generated using the scripts in the acl/ directory.")
	serverPrivateKeyPath = flag.String("server-private-key", "", "Path to server private key generated using the scripts in the acl/ directory.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	glog.Infof("starting gRPC service at :%d", *grpcPort)
	config := cloudconfig.DefaultParams()

	var grpcOpts []grpc.ServerOption
	if !*insecure {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials for gRPC: %w", err)
		}
		grpcOpts = append(grpcOpts, grpc.Creds(tlsCredentials))
		glog.Infof("starting up in secure mode using cert at %q, private key at %q", *serverCertPath, *serverPrivateKeyPath)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		s := grpc.NewServer(grpcOpts...)
		svc, err := stateservice.New(ctx, config)
		if err != nil {
			return fmt.Errorf("error initializing state service: %w", err)
		}
		pb.RegisterStateServiceServer(s, svc)
		reflection.Register(s)

		if err := grpcserverutil.ServeUntilCancelled(ctx, s, lis); err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	})
	return eg.Wait()
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(*serverCertPath, *serverPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error loading key pair: %w", err)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    acls.ClientCACertPool,
	}), nil
}
