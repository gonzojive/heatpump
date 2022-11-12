// Program cloud-listener subscribes to pub/sub thermostat commands and
// dispatches them to the fancoil service.
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/acls/deviceauth"
	"github.com/gonzojive/heatpump/cloud/queue/queueclient"
	"github.com/gonzojive/heatpump/util/grpcserverutil"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	qpb "github.com/gonzojive/heatpump/proto/command_queue"
	cpb "github.com/gonzojive/heatpump/proto/controller"
	fcpb "github.com/gonzojive/heatpump/proto/fancoil"
)

var (
	insecure              = flag.Bool("command-queue-insecure", false, "If true, don't use TLS connection to command queue server.")
	fancoilServiceAddress = flag.String("fancoil-service-addr", "192.168.86.24:8083", "Remote address of fancoil service.")
	queueAddr             = flag.String("command-queue-addr", "localhost:8083", "Remote address of CommandQueueService.")
	stateServiceAddr      = grpcserverutil.AddressFlag(
		"state-service-addr", "localhost:8083", "Remote address of StateService.",
		func(ctx context.Context, conn *grpc.ClientConn) (cpb.StateServiceClient, error) {
			return cpb.NewStateServiceClient(conn), nil
		})

	authClientFlags = deviceauth.RegisterFlagsWithPrefix(flag.CommandLine, "")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	creds, err := deviceauth.LoadTLSCredentials(authClientFlags)
	if err != nil {
		return err
	}
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	cmdClient, err := dialCommandQueue(ctx, authClientFlags)
	if err != nil {
		return err
	}

	fancoilClient, err := dialFanCoilService(ctx)
	if err != nil {
		return err
	}

	stateClient, err := stateServiceAddr.Dial(ctx, dialOpts...)
	if err != nil {
		return err
	}

	authClient, err := deviceauth.NewClient(ctx, authClientFlags)
	if err != nil {
		return fmt.Errorf("error creating deviceauth.Client: %w", err)
	}

	c := &listener{authClient, stateClient, fancoilClient}

	{
		ctx, err := authClient.AddDeviceAuthTokenToContext(ctx)
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
			if req := parsed.GetSetStateRequest(); req != nil {
				if err := c.handleSetStateRequest(ctx, parsed.GetSetStateRequest()); err != nil {
					cmd.Nack()
					glog.Errorf("failed to handle set_state_request: %v", err)
					return
				}
				cmd.Ack()
			}
		}); err != nil {
			return err
		}
	}
	return nil
}

type listener struct {
	authClient    *deviceauth.Client
	stateClient   cpb.StateServiceClient
	fancoilClient fcpb.FanCoilServiceClient
}

func (l *listener) handleSetStateRequest(ctx context.Context, req *fcpb.SetStateRequest) error {
	parsed := &cpb.Command{}
	glog.Infof("setting target temperature to %vC", parsed.GetSetStateRequest().GetHeatingTargetTemperature().GetDegreesCelcius())
	if _, err := l.fancoilClient.SetState(ctx, req); err != nil {
		return fmt.Errorf("command failed to execute: %w", err)
	}
	stateResp, err := l.fancoilClient.GetState(ctx, &fcpb.GetStateRequest{
		FancoilName: req.GetFancoilName(),
	})
	if err != nil {
		return fmt.Errorf("failed to get fan coil state: %w", err)
	}
	if _, err := l.stateClient.SetDeviceState(ctx, &cpb.SetDeviceStateRequest{
		State: &cpb.DeviceState{
			Name:         req.GetFancoilName(),
			FancoilState: stateResp.GetState(),
		},
	}); err != nil {
		return fmt.Errorf("error updating state in cloud service: %w", err)
	}
	return nil
}

func dialCommandQueue(ctx context.Context, f *deviceauth.Flags) (*queueclient.Client, error) {
	var opts []grpc.DialOption
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		creds, err := deviceauth.LoadTLSCredentials(f)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	conn, err := grpc.DialContext(ctx, *queueAddr, opts...)
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
