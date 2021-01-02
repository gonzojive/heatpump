// Program cx34install installs a systemd service for cx34collector according to
// https://www.raspberrypi.org/documentation/linux/usage/systemd.md.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"github.com/gonzojive/heatpump/linuxserial"
)

const systemdDir = "/etc/systemd/system"

var (
	rs484TTYModbus = flag.String("modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	dbDir          = flag.String("db-dir", "/home/pi/db/cx34db", "Path to a directory where badger database should be stored.")
	grpcPort       = flag.Int("grpc-port", 8082, "Port used to serve historical database values over GRPC.")
	binPath        = flag.String("collector-bin", "/home/pi/go/bin/cx34collector", "Path to collector binary")
	scriptName     = flag.String("script-name", "cx34collector.service", "Name of the systemd file to place in "+systemdDir)
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	dev, err := linuxserial.FromPath(*rs484TTYModbus)

	dev, err = dev.Canonicalize()
	if err != nil {
		return err
	}
	executable, err := canonicalizeBinPath()
	if err != nil {
		return err
	}
	glog.Infof("using executable path %s - remember to install the latest version using go install cmd/cx34collector/*.go", executable)

	scriptContents := genSystemDScript(*dbDir, dev.Path(), executable)
	scriptPath := filepath.Join(systemdDir, *scriptName)
	glog.Infof("Copying service script to %s:\n\n%s", scriptPath, scriptContents)

	if err := ioutil.WriteFile(scriptPath, []byte(scriptContents), 0664); err != nil {
		return fmt.Errorf("error writint to output file %q: %w", scriptPath, err)
	}

	if err := exec.CommandContext(ctx, "systemctl", "enable", *scriptName).Run(); err != nil {
		return fmt.Errorf("error executing `systemctl start %s`: %w", *scriptName, err)
	}
	glog.Infof("Service %s will start at each boot", *scriptName)

	if err := exec.CommandContext(ctx, "systemctl", "restart", *scriptName).Run(); err != nil {
		return fmt.Errorf("error executing `systemctl restart %s`: %w", *scriptName, err)
	}
	glog.Infof("Service %s started using `systemctrl restart %s`", *scriptName, *scriptName)

	fmt.Printf(`The service has been started and enabled (will start at each boot). To restart, execute
	
	sudo systemctl start %s

Stop the service using the following command:

	sudo systemctl stop %s

View logs of the service using the following command:

	journalctl -u %s
`, *scriptName, *scriptName, *scriptName)

	return nil
}

func canonicalizeBinPath() (string, error) {
	_, err := os.Stat(*binPath)
	if err != nil {
		return "", fmt.Errorf("error verifying collector-bin is a valid file: %w", err)
	}
	return filepath.Abs(*binPath)
}

func genSystemDScript(dbDir, modbusTTYDevice, collectorBinPath string) string {
	return fmt.Sprintf(`# systemd file for CX34 heat pump chilelr generated on %s
# See https://www.raspberrypi.org/documentation/linux/usage/systemd.md for
# instructions about using systemd with Raspberry Pi.
#
# Project home page: https://github.com/gonzojive/heatpump
[Unit]
Description=Chiltrix CX34 heat pump controller and sensor recorder
After=network.target

[Service]
ExecStart=%s --db-dir %q --modbus-device %q --alsologtostderr
WorkingDirectory=/home/pi
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
`, time.Now().Format(time.RFC3339), collectorBinPath, dbDir, modbusTTYDevice)
}
