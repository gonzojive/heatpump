// Program cx34runs a dashboard that displays information about a CX34 heat pump.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/golang/glog"
	"github.com/martinlindhe/unit"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/dashboard"
	"github.com/gonzojive/heatpump/fancoil"
	"github.com/gonzojive/heatpump/grpcspec"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"go.uber.org/fx"

	fancoilpb "github.com/gonzojive/heatpump/proto/fancoil"
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

	fancoilServiceParams = grpcspec.ClientParamsFlag{
		Value: grpcspec.NewClientParamsBuilder().
			SetAddr("localhost:8083").
			SetInsecure(true).
			Build(),
	}
)

func init() {
	flag.Var(&readWriteServiceParams, "read-write-grpc-service", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true) for chiltrix.ReadWriteService")
	flag.Var(&fancoilServiceParams, "fan-coil-grpc-service", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true) for FanCoilService")
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
			fx.Provide(newFancoilServiceClient),
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

func newFancoilServiceClient(lc fx.Lifecycle) (fancoilpb.FanCoilServiceClient, error) {
	client, conn, err := grpcspec.NewClient(fancoilServiceParams.Value, fancoilpb.NewFanCoilServiceClient)
	if err != nil {
		return nil, fmt.Errorf("error creating FanCoilServiceClient: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return client, nil
}

func NewHTTPServer(lc fx.Lifecycle, rwServiceClient chiltrix.ReadWriteServiceClient, fancoilClient fancoilpb.FanCoilServiceClient) *http.Server {
	mux := http.NewServeMux()
	// Add your handler for /index.md here
	mux.HandleFunc("/index.md", func(w http.ResponseWriter, r *http.Request) {
		markdown, err := generateMarkdown(r.Context(), rwServiceClient, fancoilClient)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate markdown: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/markdown")
		fmt.Fprint(w, markdown)
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

func generateMarkdown(
	ctx context.Context,
	rwServiceClient chiltrix.ReadWriteServiceClient,
	fancoilClient fancoilpb.FanCoilServiceClient) (string, error) {
	resp, err := rwServiceClient.GetState(ctx, &chiltrix.GetStateRequest{})
	if err != nil {
		return "", fmt.Errorf("error getting state from heat pump: %w", err)
	}
	state, err := cx34.StateFromProto(resp)
	if err != nil {
		return "", fmt.Errorf("error parsing state from heat pump: %w", err)
	}
	fancoilReportMarkdown, err := fancoilsReport(ctx, fancoilClient)
	if err != nil {
		return "", fmt.Errorf("error parsing state from fancoil service: %w", err)
	}
	return fmt.Sprintf(`# HVAC System Dashboard

## Fan coil state

%s

%s`, fancoilReportMarkdown, state.Report(false, nil)), nil
}

func fancoilsReport(
	ctx context.Context,
	fancoilClient fancoilpb.FanCoilServiceClient) (string, error) {

	var states []*fancoil.State
	var reports []string

	{
		resp, err := fancoilClient.GetState(ctx, &fancoilpb.GetStateRequest{})

		if err != nil {
			return "", fmt.Errorf("error getting state from fancoils: %w", err)
		}
		state, err := fancoil.StateFromSnapshotProto(resp.State.SnapshotTime.AsTime(), resp.GetRawSnapshot())
		if err != nil {
			return "", fmt.Errorf("error parsing state from fancoils: %w", err)
		}
		states = append(states, state)
	}

	for _, state := range states {
		reports = append(reports, fancoilReport(state))
	}

	return strings.Join(reports, "\n\n"), nil
}

func fancoilReport(state *fancoil.State) string {
	stateProto := state.ResponseProto()
	return fmt.Sprintf(
		`- Power status: %s
- Fan setting: %s (preference) %s (current)
- Fan speed: %d RPM
- Heating/cooling target temp: %s / %s
- Room temp: %s
- Hydronic coil temperature: %s`,
		stateProto.GetState().GetPowerStatus(),
		stateProto.GetState().GetPreferenceFanSetting(),
		stateProto.GetState().GetCurrentFanSetting(),
		stateProto.GetState().GetFanSpeed().GetRpm(),
		formatTemp(parseTempProto(stateProto.GetState().GetHeatingTargetTemperature())),
		formatTemp(parseTempProto(stateProto.GetState().GetCoolingTargetTemperature())),
		formatTemp(parseTempProto(stateProto.GetState().GetRoomTemperature())),
		formatTemp(parseTempProto(stateProto.GetState().GetCoilTemperature())),
	)
}

func formatTemp(temp unit.Temperature) string {
	return fmt.Sprintf("%.1f°C (%.1f°F)", temp.Celsius(), temp.Fahrenheit())
}

func parseTempProto(temp *fancoilpb.Temperature) unit.Temperature {
	return unit.FromCelsius(float64(temp.GetDegreesCelcius()))
}
