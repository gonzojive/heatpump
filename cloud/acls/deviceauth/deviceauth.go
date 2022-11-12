// Package deviceauth is a library for obtaining credentials on an IoT device
// for authenticating with some of the cloud services like StateService and
// CommandQueueService.
package deviceauth

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/gonzojive/heatpump/proto/controller"
	"github.com/gonzojive/heatpump/util/grpcserverutil"
)

// DeviceAccessTokenMetadataKey is the gRPC metadata key used to provide a
// device access token identity secret to the server.
const DeviceAccessTokenMetadataKey = "heatpump-device-access-token"

type Flags struct {
	// For now we store a secret token in a file and refresh it. However, we
	// should move to using a private key to obtain a token so that if the token
	// is stolen the device is still the only entity able to authenticate after
	// the stolen token expires.
	fixmeTokenPath string

	authServiceAddr string
}

// RegisterFlagsWithPrefix registers a set of flags for obtaining configuration
// information.
func RegisterFlagsWithPrefix(fs *flag.FlagSet, prefix string) *Flags {
	f := &Flags{}
	fs.StringVar(&f.fixmeTokenPath, prefix+"secret-token-path", "", "Path to a file contain a secret token used for authenticating.")
	fs.StringVar(&f.authServiceAddr, prefix+"auth-service-addr", "", "Address of running gRPC AuthSerivce.")
	return f
}

type Client struct {
	// wantGracefulStop is closed when GracefulStop iscalled.
	wantGracefulStop chan struct{}
	gracefulStopOnce sync.Once

	activeDeviceToken string
	tokenExp          time.Time
	// protobuf client
	c pb.AuthServiceClient
}

func NewClient(ctx context.Context, flags *Flags) (*Client, error) {
	if flags.fixmeTokenPath == "" {
		return nil, fmt.Errorf("must specify secret-token-path flag")
	}
	tokenBytes, err := ioutil.ReadFile(flags.fixmeTokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read auth secret: %w", err)
	}

	pbClient, err := grpcserverutil.DialSecure(ctx, flags.authServiceAddr, func(conn *grpc.ClientConn) (pb.AuthServiceClient, error) {
		return pb.NewAuthServiceClient(conn), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth service: %w", err)
	}

	return &Client{
		make(chan struct{}),
		sync.Once{},
		strings.TrimSpace(string(tokenBytes)),
		time.Now().Add(time.Hour * 24 * 700),
		pbClient,
	}, nil
}

// AddDeviceAuthTokenToContext returns a new context with a device auth token in the metadata.
func (c *Client) AddDeviceAuthTokenToContext(ctx context.Context) (context.Context, error) {
	if c.tokenExp.Before(time.Now().Add(time.Minute * 60)) {
		// c.c.ExtendToken(ctx, &pb.ExtendTokenRequest{
		// 	Token: c.activeDeviceToken,
		// }, opts ...grpc.CallOption)
		return nil, fmt.Errorf("device token needs to be refresehd (unsupported)")
	}
	return metadata.AppendToOutgoingContext(ctx, DeviceAccessTokenMetadataKey, c.activeDeviceToken), nil
}

// Run keeps a device token active and sends events to the callback.
// func (c *Client) Run(ctx context.Context) error {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	ticker := time.NewTicker(time.Hour)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <- c.wantGracefulStop:
// 			return nil
// 		case <- ctx.Done():
// 			return ctx.Err()
// 		case <- ticker.C:

// 		}
// 		resp, err := c.c.ExtendToken(ctx, in *pb.ExtendTokenRequest, opts ...grpc.CallOption)
// 	}
// }

// GracefulStop attempts to gracefully shut down the server.
func (c *Client) GracefulStop() {
	c.gracefulStopOnce.Do(func() {
		close(c.wantGracefulStop)
	})
}

// LoadTLSCredentials returns the standard transport credentials for connecting
// to a TLS server.
func LoadTLSCredentials(flags *Flags) (credentials.TransportCredentials, error) {
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

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

/*

TODO(reddaly): Move the following into a library for working with mutual TLS
using private client and server certs.

	clientCertPath               = flag.String("client-cert", "", "Path to client cert generated using the scripts in the acl/ directory.")
	clientPrivateKeyPath         = flag.String("client-private-key", "", "Path to client private key generated using the scripts in the acl/ directory.")
	cloudServerAuthorityCertPath = flag.String("server-ca-cert", "", "Path to cert that signed the server's TLS cert.")

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate.
	certPool, err := func() (*x509.CertPool, error) {
		if *cloudServerAuthorityCertPath == "" {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				return nil, fmt.Errorf("failed to load system cert pool: %w", err)
			}
			return certPool, nil
		}
		pemServerCA, err := ioutil.ReadFile(*cloudServerAuthorityCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load cert of CA who signed server's certificate: %w", err)
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			return nil, fmt.Errorf("failed to add server CA's certificate")
		}
		return certPool, nil

	}()

	if err != nil {
		return nil, err
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(*clientCertPath, *clientPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert and/or private key: %w", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

*/
