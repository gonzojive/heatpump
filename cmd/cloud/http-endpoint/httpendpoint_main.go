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
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/httpendpoint"
)

var (
	port = flag.Int("port", 42598, "Local port to use for HTTP server.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {

	serveMux := http.DefaultServeMux

	_, err := httpendpoint.NewServer(ctx, serveMux, cloudconfig.DefaultParams())
	if err != nil {
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
