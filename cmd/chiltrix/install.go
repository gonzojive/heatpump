package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"github.com/gonzojive/heatpump/linuxserial"
	"github.com/spf13/cobra"
)

const systemdDir = "/etc/systemd/system"

var (
	instRs484TTYModbus string
	instDbDir          string
	instGrpcPort       int
	instBinPath        string
	instScriptName     string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a systemd service for cx34 controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		dev, err := linuxserial.FromPath(instRs484TTYModbus)
		if err == nil {
			if devC, errC := dev.Canonicalize(); errC == nil {
				instRs484TTYModbus = devC.Path()
			}
		}

		executable, err := canonicalizeBinPath(instBinPath)
		if err != nil {
			return err
		}
		glog.Infof("using executable path %s - remember to build via bazel", executable)

		scriptContents := genSystemDScript(instDbDir, instRs484TTYModbus, executable)
		scriptPath := filepath.Join(systemdDir, instScriptName)
		glog.Infof("Copying service script to %s:\n\n%s", scriptPath, scriptContents)

		if err := ioutil.WriteFile(scriptPath, []byte(scriptContents), 0664); err != nil {
			return fmt.Errorf("error writing to output file %q: %w", scriptPath, err)
		}

		if err := exec.CommandContext(ctx, "systemctl", "enable", instScriptName).Run(); err != nil {
			return fmt.Errorf("error executing `systemctl enable %s`: %w", instScriptName, err)
		}
		glog.Infof("Service %s will start at each boot", instScriptName)

		if err := exec.CommandContext(ctx, "systemctl", "restart", instScriptName).Run(); err != nil {
			return fmt.Errorf("error executing `systemctl restart %s`: %w", instScriptName, err)
		}
		glog.Infof("Service %s started using `systemctl restart %s`", instScriptName, instScriptName)

		fmt.Printf(`The service has been started and enabled (will start at each boot). To restart, execute
	
	sudo systemctl start %s

Stop the service using the following command:

	sudo systemctl stop %s

View logs of the service using the following command:

	journalctl -u %s
`, instScriptName, instScriptName, instScriptName)

		return nil
	},
}

func init() {
	installCmd.Flags().StringVar(&instRs484TTYModbus, "modbus-device", "/dev/ttyUSB0", "Path to USB-to-RS485 device connected to modbus.")
	installCmd.Flags().StringVar(&instDbDir, "db-dir", "/home/pi/db/cx34db", "Path to a directory where badger database should be stored.")
	installCmd.Flags().IntVar(&instGrpcPort, "grpc-port", 8082, "Port used to serve historical database values over GRPC.")
	installCmd.Flags().StringVar(&instBinPath, "collector-bin", "/home/pi/bin/chiltrix", "Path to chiltrix main binary")
	installCmd.Flags().StringVar(&instScriptName, "script-name", "chiltrix.service", "Name of the systemd file to place in "+systemdDir)
	rootCmd.AddCommand(installCmd)
}

func canonicalizeBinPath(binPath string) (string, error) {
	_, err := os.Stat(binPath)
	if err != nil {
		return "", fmt.Errorf("error verifying collector-bin is a valid file: %w", err)
	}
	return filepath.Abs(binPath)
}

func genSystemDScript(dbDir, modbusTTYDevice, collectorBinPath string) string {
	return fmt.Sprintf(`# systemd file for CX34 heat pump chiller generated on %s
# See https://www.raspberrypi.org/documentation/linux/usage/systemd.md for
# instructions about using systemd with Raspberry Pi.
#
# Project home page: https://github.com/gonzojive/heatpump
[Unit]
Description=Chiltrix CX34 heat pump unified controller
After=network.target

[Service]
ExecStart=%s start-readwrite-service --modbus-device %q --alsologtostderr
WorkingDirectory=/home/pi
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
`, time.Now().Format(time.RFC3339), collectorBinPath, modbusTTYDevice)
}
