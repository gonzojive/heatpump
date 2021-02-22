// Program cx34runs a dashboard that displays information about a CX34 heat pump.
package main

import (
	"context"
	"flag"
	"time"

	"github.com/golang/glog"

	"github.com/gonzojive/heatpump/dashboard"
)

var (
	httpPort        = flag.Int("port", 8081, "HTTP server port.")
	historianAddr   = flag.String("collector", ":8082", "Address of the cx34collector gRPC service. May be on a remote machine.")
	dashboardWindow = flag.Duration("dashboard_window", time.Hour*6, "Amount of data to display on the dashboard.")
)

func main() {
	flag.Parse()
	if err := dashboard.Run(context.Background(), *historianAddr, *httpPort, *dashboardWindow); err != nil {
		glog.Exitf("%v", err)
	}
}
