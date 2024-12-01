// Program cx34runs a dashboard that displays information about a CX34 heat pump.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/golang/glog"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/dashboard"
	"github.com/gonzojive/heatpump/grpcspec"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"go.uber.org/fx"
)

var (
	httpPort         = flag.Int("port", 8081, "HTTP server port.")
	historianAddr    = flag.String("collector", ":8082", "Address of the cx34collector gRPC service. May be on a remote machine.")
	dashboardVersion = flag.String("version", "2", "Dashboard application version to launch. Must be '1' or '2'.")

	readWriteServiceParams = grpcspec.ClientParamsFlag{
		Value: grpcspec.NewClientParamsBuilder().
			SetAddr("localhost:8084").
			SetInsecure(true).
			Build(),
	}
)

func init() {
	flag.Var(&readWriteServiceParams, "read-write-grpc-service", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true)")
}

func main() {
	flag.Parse()
	switch *dashboardVersion {
	case "1":
		glog.Infof("starting up heat pump dashboard at http://localhost:%d", *httpPort)
		if err := dashboard.Run(context.Background(), *historianAddr, *httpPort); err != nil {
			glog.Exitf("%v", err)
		}
	case "2":
		fx.New(
			fx.Provide(NewHTTPServer),
			fx.Provide(newReadWriteServiceClient),
			fx.Invoke(func(*http.Server) {}),
		).Run()
	default:
		glog.Exitf("Invalid version %q", *dashboardVersion)
	}

}

func newReadWriteServiceClient(lc fx.Lifecycle) (chiltrix.ReadWriteServiceClient, error) {
	client, conn, err := grpcspec.NewClient(readWriteServiceParams.Value, chiltrix.NewReadWriteServiceClient)
	if err != nil {
		return nil, fmt.Errorf("error creating ReadWriteServiceClient: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return client, nil
}

func NewHTTPServer(lc fx.Lifecycle, rwServiceClient chiltrix.ReadWriteServiceClient) *http.Server {
	mux := http.NewServeMux()
	// Add your handler for /index.md here
	mux.HandleFunc("/index.md", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown")
		fmt.Fprint(w, "# Chiltrix heatpump dashboard\n\n")
		resp, err := rwServiceClient.GetState(r.Context(), &chiltrix.GetStateRequest{})
		if err != nil {
			fmt.Fprintf(w, "Error getting state from heat pump: %v", err)
			return
		}
		state, err := cx34.StateFromProto(resp)
		if err != nil {
			fmt.Fprintf(w, "Error parsing state from heat pump: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", state.Report(false, nil))
	})
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		// Get the runfiles manifest
		rfs, err := runfiles.New()
		if err != nil {
			http.Error(w, "Failed to get runfiles", http.StatusInternalServerError)
			return
		}

		// Load the runfile using the manifest
		indexHTMLPath, err := rfs.Rlocation("github-gonzojive-heatpump/cmd/cx34dash/cx34dash.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load cx34dash.html: %v", err), http.StatusInternalServerError)
			return
		}
		indexHTML, err := os.ReadFile(indexHTMLPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load cx34dash.html: %v", err), http.StatusInternalServerError)
			return
		}

		// Output the content
		w.Write(indexHTML)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: mux,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
