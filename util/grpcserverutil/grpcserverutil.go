// Package grpcserverutil contains utilities for running gRPC servers.
//
// TODO(reddaly): Rename package gcputil.
package grpcserverutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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

// DialSecure dials a gRPC service securely using the system cert pool and calls
// the factory method to construct a client.
func DialSecure[T any](ctx context.Context, addr string, factory func(conn *grpc.ClientConn) (T, error)) (T, error) {
	zero := func() T {
		var t T
		return t
	}
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return zero(), fmt.Errorf("failed to load system cert pool: %w", err)
	}

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})))
	if err != nil {
		return zero(), fmt.Errorf("failed to connect to FanCoilService: %w", err)
	}
	return factory(conn)
}

// SystemCertPoolTransportCredentials returns the system TransportCredentials.
func SystemCertPoolTransportCredentials() (credentials.TransportCredentials, error) {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system cert pool: %w", err)
	}

	return credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	}), nil
}

// Address is primarily used for defining flags that can be used to instantiate
// gRPC clients.
type Address[ClientT any] struct {
	addr    string
	factory func(ctx context.Context, conn *grpc.ClientConn) (ClientT, error)
}

// Dial creates a new connection to the provided address and returns a client
// that uses that connection.
func (a *Address[ClientT]) Dial(ctx context.Context, opts ...grpc.DialOption) (ClientT, error) {
	conn, err := grpc.DialContext(ctx, a.addr, opts...)
	if err != nil {
		var zero ClientT
		return zero, fmt.Errorf("error dialing %q: %w", a.addr, err)
	}
	return a.factory(ctx, conn)
}

// AddressFlag defines a new flag of type Address[ClientT].
func AddressFlag[ClientT any](
	name string,
	defaultVal string,
	usage string,
	factory func(ctx context.Context, conn *grpc.ClientConn) (ClientT, error)) *Address[ClientT] {
	return RegisterAddressFlag(flag.CommandLine, name, defaultVal, usage, factory)
}

// RegisterAddressFlag defines a new flag of type Address[ClientT].
func RegisterAddressFlag[ClientT any](
	fs *flag.FlagSet,
	name string,
	defaultVal string,
	usage string,
	factory func(ctx context.Context, conn *grpc.ClientConn) (ClientT, error)) *Address[ClientT] {

	f := &Address[ClientT]{defaultVal, factory}
	fs.StringVar(&f.addr, name, defaultVal, usage)
	return f
}

// TransportCredentialsSpec specifies TransportCredentials.
type TransportCredentialsSpec string

const (
	// TLS specifies a TLS connection secured using the system TLS cert pool.
	TLS TransportCredentialsSpec = "tls"

	// Insecure specifies a connection that doesn't use any encryption.
	Insecure TransportCredentialsSpec = "insecure"
)

func (s TransportCredentialsSpec) DialOptions() ([]grpc.DialOption, error) {
	switch s {
	case TLS:
		creds, err := SystemCertPoolTransportCredentials()
		if err != nil {
			return nil, fmt.Errorf("error loading system cert pool for gRPC TLS connections: %w", err)
		}
		return []grpc.DialOption{grpc.WithTransportCredentials(creds)}, nil
	case Insecure:
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	default:
		return nil, fmt.Errorf("invalid TransportCredentialsSpec %q", s)
	}
}

type AddressSpec struct {
	Address string
}

// SystemTLSCredentials returns the standard transport credentials for connecting
// to a TLS server.
func SystemTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate.
	certPool, err := func() (*x509.CertPool, error) {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to load system cert pool: %w", err)
		}
		return certPool, nil
	}()

	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	}), nil
}
