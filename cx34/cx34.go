// Package cx34 provides a client for working with the Chiltrix CX34 heat pump.
package cx34

import (
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"github.com/golang/glog"
	"github.com/howeyc/crc16"
)

// Parameters from https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf
const (
	baudRate = 9600
	parity   = "N"
	stopBits = 1
	dataBits = 8
	slaveID  = 1
)

type Mode string

const (
	Modbus   Mode = "modbus"
	CX34Text Mode = "cx34text"
)

// Params configures the connection to the Chiltrix.
type Params struct {
	// The /dev/ttyX device shown by dmesg for the RS-485 connection to the heat pump.
	TTYDevice string
	LogWriter io.Writer
	Mode      Mode
}

// Client is used to communicate with the Chiltrix CX34 heat pump.
type Client struct {
	c modbus.Client
}

// Connect connects a new client to the heat pump or returns an error.
func Connect(p *Params) (*Client, error) {
	if p.Mode != Modbus && p.Mode != CX34Text {
		return nil, fmt.Errorf("Invalid mode %q", p.Mode)
	}
	// Modbus RTU/ASCII
	handler := modbus.NewRTUClientHandler(p.TTYDevice)
	handler.BaudRate = baudRate
	handler.DataBits = dataBits
	handler.Parity = parity
	handler.StopBits = stopBits
	handler.SlaveId = slaveID
	handler.Timeout = 5 * time.Second
	handler.RS485 = serial.RS485Config{
		Enabled:            false,
		DelayRtsAfterSend:  0,
		DelayRtsBeforeSend: 0,
		RtsHighAfterSend:   false,
		RtsHighDuringSend:  false,
		RxDuringTx:         true,
	}

	if p.Mode == CX34Text {
		port, err := serial.Open(&handler.Config)
		if err != nil {
			return nil, fmt.Errorf("serial.Open failure: %w", err)
		}
		_, err = io.Copy(p.LogWriter, port)
		if err != nil {
			return nil, fmt.Errorf("io.Copy failure: %w", err)
		}
		return nil, nil
	}

	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		return nil, fmt.Errorf("Connect failed: %w", err)
	}

	client := modbus.NewClient(handler)
	results, err := client.ReadCoils(0, 1)
	if err != nil {
		return nil, fmt.Errorf("ReadCoils() failed: %w", err)
	}
	glog.Infof("got results %v", results)
	return &Client{client}, nil
}

type Logger struct {
	debug, raw io.Writer
}

func NewLogger(debug, raw io.Writer) *Logger {
	return &Logger{debug, raw}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	t := time.Now()
	//n, err = l.data.Write(p)
	{
		pdu, err := decodeFrame(p)
		msg := fmt.Sprintf("error decoding frame: %v", err)
		if err == nil {
			msg = fmt.Sprintf("decoded frame with code %v", pdu.FunctionCode)
		}
		glog.Infof("got %q: %s", string(p), msg)
	}
	if _, err := fmt.Fprintf(l.debug, "%d %q // %d bytes; %s\n", t.UnixNano(), hex.EncodeToString(p), len(p), t); err != nil {
		return 0, err
	}
	if l, err := l.raw.Write(p); err != nil {
		return 0, fmt.Errorf("raw output error: %w", err)
	} else if l != len(p) {
		return 0, fmt.Errorf("raw output Write() returned %d, want %d", l, len(p))
	}
	return len(p), nil
}

func decodeFrame(adu []byte) (*modbus.ProtocolDataUnit, error) {
	if len(adu) < 4 {
		return nil, fmt.Errorf("modbus: argument cannot possibly be a legitimate Application Data Unit: size is too small (%d bytes)", len(adu))
	}
	length := len(adu)
	// Calculate checksum
	crcCalculated := ^crc16.ChecksumIBM(adu[0 : length-2]) //crc16.Update(65535, crc16.IBMTable, adu[0:length-2])
	crcPacket := uint16(adu[length-1])<<8 | uint16(adu[length-2])
	if crcPacket != crcCalculated {
		return nil, fmt.Errorf("modbus: response crc %d does not match expected %d in %d-byte package", crcCalculated, crcPacket, length)
	}
	// Function code & data
	return &modbus.ProtocolDataUnit{
		FunctionCode: adu[1],
		Data:         adu[2 : length-2],
	}, nil
}
