// Program cx34runs a dashboard that displays information about a CX34 heat pump.
package main

import (
	"context"
	"flag"

	"github.com/golang/glog"

	"github.com/gonzojive/heatpump/dashboard"
)

var (
	httpPort      = flag.Int("port", 8081, "HTTP server port.")
	historianAddr = flag.String("collector", ":8082", "Address of the cx34collector gRPC service. May be on a remote machine.")
)

func main() {
	flag.Parse()
	glog.Infof("starting up heat pump dashboard at http://localhost:%d", *httpPort)
	if err := dashboard.Run(context.Background(), *historianAddr, *httpPort); err != nil {
		glog.Exitf("%v", err)
	}
}
