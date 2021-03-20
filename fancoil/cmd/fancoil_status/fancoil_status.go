// Program fancoil_status prints out the status of any connected fan coil.
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/fancoil"
)

var (
	rs484TTYModbus = flag.String("modbus-device", "/dev/serial/by-id/usb-1a86_USB_Serial-if00-port0", "Path to USB-to-RS485 device connected to modbus.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}
func run(ctx context.Context) error {
	_, err := fancoil.Connect(&fancoil.Params{
		TTYDevice: *rs484TTYModbus,
	})
	if err != nil {
		return fmt.Errorf("error connecting: %w", err)
	}
	return nil
}
