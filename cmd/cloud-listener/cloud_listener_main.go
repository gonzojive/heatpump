// Program cloud-listener subscribes to pub/sub thermostat commands and
// dispatches them to the fancoil service.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/queue/queueclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"

	qpb "github.com/gonzojive/heatpump/proto/command_queue"
	cpb "github.com/gonzojive/heatpump/proto/controller"
	fcpb "github.com/gonzojive/heatpump/proto/fancoil"
)

var (
	fancoilServiceAddress = flag.String("fancoil-service-addr", "192.168.86.24:8083", "Remote address of fancoil service.")
	queueAddr             = flag.String("command-queue-addr", "localhost:8083", "Remote address of CommandQueueService.")

	clientCertPath               = flag.String("client-cert", "", "Path to client cert generated using the scripts in the acl/ directory.")
	clientPrivateKeyPath         = flag.String("client-private-key", "", "Path to client private key generated using the scripts in the acl/ directory.")
	cloudServerAuthorityCertPath = flag.String("server-ca-cert", "", "Path to cert that signed the server's TLS cert.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	cmdClient, err := dialCommandQueue(ctx)
	if err != nil {
		return err
	}

	fancoilClient, err := dialFanCoilService(ctx)
	if err != nil {
		return err
	}

	if err := cmdClient.Listen(ctx, func(cmd *queueclient.Command) {
		parsed := &cpb.Command{}
		if err := proto.Unmarshal(cmd.Payload(), parsed); err != nil {
			glog.Errorf("error processing command: invalid payload: %v", err)
			// Invalid now, invalid forever. Ack to avoid a backlog of invalid
			// messages.
			cmd.Ack()
			return
		}
		if setStateReq := parsed.GetSetStateRequest(); setStateReq != nil {
			glog.Infof("setting target temperature to %vC", parsed.GetSetStateRequest().GetHeatingTargetTemperature().GetDegreesCelcius())
			if _, err := fancoilClient.SetState(ctx, setStateReq); err != nil {
				glog.Warningf("command failed to execute: %w", err)
			}
		}
		cmd.Ack()
	}); err != nil {
		return err
	}
	return nil
}

func dialCommandQueue(ctx context.Context) (*queueclient.Client, error) {
	creds, err := loadTLSCredentials()
	if err != nil {
		return nil, err
	}
	conn, err := grpc.DialContext(ctx, *queueAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FanCoilService: %w", err)
	}
	return queueclient.New(qpb.NewCommandQueueServiceClient(conn)), nil
}

func dialFanCoilService(ctx context.Context) (fcpb.FanCoilServiceClient, error) {
	conn, err := grpc.DialContext(ctx, *fancoilServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FanCoilService: %w", err)
	}
	return fcpb.NewFanCoilServiceClient(conn), nil
}

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
