// Program http-endpoint is an OAuth 2 server for operating a
// home's fan coil units using Google Asistant.
//
// See https://developers.google.com/assistant/smarthome/overview.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/acls/server2serverauth"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/httpendpoint"
	"github.com/gonzojive/heatpump/util/grpcserverutil"
	"google.golang.org/grpc"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

var (
	port = flag.Int("port", 42598, "Local port to use for HTTP server.")

	stateServiceAddr = grpcserverutil.AddressFlag(
		"state-service-addr", "localhost:8089", "Remote address of StateService.",
		func(ctx context.Context, conn *grpc.ClientConn) (cpb.StateServiceClient, error) {
			return cpb.NewStateServiceClient(conn), nil
		})

	insecureS2S = flag.Bool("insecure-s2s", false, "If true, disable server-to-server auth.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	serveMux := http.DefaultServeMux

	var ssOpts []grpc.DialOption

	{
	}
	if *insecureS2S {
		ssOpts = append(ssOpts, grpc.WithInsecure())
	} else {
		creds, err := grpcserverutil.SystemTLSCredentials()
		if err != nil {
			return err
		}
		ssOpts = append(ssOpts, grpc.WithTransportCredentials(creds))

		sender, err := server2serverauth.NewSender(ctx, "stateservice")
		if err != nil {
			return err
		}
		ssOpts = append(ssOpts, grpc.WithPerRPCCredentials(sender.PerRPCCredentials()))
	}

	ssc, err := stateServiceAddr.Dial(ctx, ssOpts...) //, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("error dialing StateService: %w", err)
	}

	if _, err := httpendpoint.NewServer(ctx, serveMux, cloudconfig.DefaultParams(), ssc); err != nil {
		return err
	}

	finalPort := *port

	if got := os.Getenv("PORT"); got != "" {
		x, err := strconv.Atoi(got)
		if err != nil {
			glog.Fatalf("bad PORT env var %q", got)
		}
		finalPort = x
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", finalPort), nil)
}
