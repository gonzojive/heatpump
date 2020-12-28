// Package linuxserial provides utilitis for working with serial devices in
// Linux. This was written to list and canonicalize RS-485 USB devices used to
// connect to a heat pump using modbus RTU.
package linuxserial

import (
	"fmt"
	"path/filepath"
)

const byIDDir = "/dev/serial/by-id"

// Device is a pointer to a serial device.
type Device struct {
	// devPath is the path to the device "file." Example value: /dev/ttyUSB0 or
	// /dev/serial/by-id/usb-FTDI_FT232R_USB_UART_XYZ-if00-port0
	devPath string
}

// FromPath returns a serial device based on a file path like "/dev/ttyUSB0" or
// /dev/serial/by-id/usb-FTDI_FT232R_USB_UART_XYZ-if00-port0.
func FromPath(path string) (Device, error) {
	return Device{path}, nil
}

// Path returns the path to the device "file." Example value: /dev/ttyUSB0 or
// /dev/serial/by-id/usb-FTDI_FT232R_USB_UART_XYZ-if00-port0.
func (d Device) Path() string {
	return d.devPath
}

// Canonicalize returns a device that points to the same underlying device but
// has a persistent pathname. The returned value typically begins with /dev/serial/by-id.
func (d Device) Canonicalize() (Device, error) {
	resolvedPath, err := filepath.EvalSymlinks(d.Path())
	if err != nil {
		return Device{}, fmt.Errorf("error resolving symlinks for %q: %w", d.Path(), err)
	}
	devs, err := List()
	if err != nil {
		return Device{}, err
	}
	for _, cand := range devs {
		candidateResolvedPath, err := filepath.EvalSymlinks(cand.Path())
		if err != nil {
			return Device{}, fmt.Errorf("error resolving symlinks for %q: %w", cand.Path(), err)
		}
		if candidateResolvedPath == resolvedPath {
			return cand, nil
		}
	}
	return Device{}, fmt.Errorf("could not find device in list: %+v", devs)
}

// List lists available serial devices.
func List() ([]Device, error) {
	paths, err := filepath.Glob(filepath.Join(byIDDir, "*"))
	if err != nil {
		return nil, fmt.Errorf("error listing %s for serial devices: %w", byIDDir, err)
	}
	var out []Device
	for _, p := range paths {
		out = append(out, Device{p})
	}
	return out, nil
}
