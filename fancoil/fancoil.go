// Package fancoil communicates with a Heat Transfer Products (and presumably
// Chiltrix) fan coil unit.
//
// The fan coil unit manual: https://htproducts.com/literature/lp-587.pdf
//
// The Chiltrix and HTP fan coil units are manufactured by the same Chinese
// company. This package should work for any fan coil unit with a ZLFP10 circuit
// board; it will not work with the older MD1001 circuit board. HTP seems to
// ship units with both circuit boards, and according to one technical support
// representative the boards are not interchangable.
//
// The modbus functionality is not officially supported by HTP. The
// documentation from the manual is incomplete but still useful.
package fancoil

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/mdtable"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"go.uber.org/multierr"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Parameters from https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf
const (
	baudRate = 9600
	parity   = "N"
	stopBits = 1
	dataBits = 8
	slaveID  = 15
)

const (
	// Valid Holding register range
	firstHoldingRegister Register = 1
	lastHoldingRegister  Register = 10
	registersPerRead              = 20
)

// Mode indicates the protocol that should be used to communicate with the CX34.
type Mode string

// Params configures the connection to the Chiltrix.
type Params struct {
	// The /dev/ttyX device shown by dmesg for the RS-485 connection to the heat pump.
	TTYDevice string
	Mode      Mode
}

// Client is used to communicate with the Chiltrix CX34 heat pump.
type Client struct {
	chiltrix.ReadWriteServiceServer
	c modbus.Client
}

// Connect connects a new client to the heat pump or returns an error.
func Connect(p *Params) (*Client, error) {
	// Modbus RTU.
	handler := modbus.NewRTUClientHandler(p.TTYDevice)
	handler.BaudRate = baudRate
	handler.DataBits = dataBits
	handler.Parity = parity
	handler.StopBits = stopBits
	handler.SlaveId = slaveID
	handler.Timeout = 10 * time.Second
	handler.RS485 = serial.RS485Config{
		Enabled:            false,
		DelayRtsAfterSend:  0,
		DelayRtsBeforeSend: 0,
		RtsHighAfterSend:   false,
		RtsHighDuringSend:  false,
		RxDuringTx:         true,
	}

	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		return nil, fmt.Errorf("Connect failed: %w", err)
	}

	client := modbus.NewClient(handler)
	c := &Client{c: client}

	if err := c.CheckConnection(); err != nil {
		return nil, err
	}
	return c, nil
}

// ReadState returns a snapshot of the state of the heat pump.
func (c *Client) ReadState() (*State, error) {
	// ReadCoils, ReadInputRegisters, and ReadDiscreteInputs are not supported.
	// However, ReadHoldingRegisters is.
	m := make(map[Register]uint16)
	for i := firstHoldingRegister; i <= lastHoldingRegister; i += registersPerRead {
		count := registersPerRead
		if i+Register(count) > lastHoldingRegister {
			count = int(lastHoldingRegister) - int(i+1)
		}
		results, err := c.c.ReadHoldingRegisters(uint16(i), uint16(count))
		if err != nil {
			return nil, fmt.Errorf("ReadHoldingRegisters() failed: %w", err)
		}
		if len(results)%2 != 0 {
			return nil, fmt.Errorf("got register data of length %d, want modulus of 2", len(results))
		}
		if len(results)/2 > count {
			return nil, fmt.Errorf("returned register data of length %d exceeds expected length %d", len(results)/2, count)
		}
		for j := 0; j < len(results)/2; j++ {
			value := uint16(results[j*2])<<8 + uint16(results[j*2+1])
			m[Register(j)+i] = value
		}
	}
	return &State{time.Now(), m}, nil
}

func (c *Client) setRegisterValue(reg, value uint16) error {
	res, err := c.c.WriteSingleRegister(reg, value)
	if err != nil {
		return fmt.Errorf("WriteSingleRegister error: %w (returned bytes %v)", err, res)
	}
	glog.Infof("set register value %d to %d; got response %v", reg, value, res)
	return nil
}

// CheckConnection attempts to connect to the heat pump and returns an error if the connection fails.
func (c *Client) CheckConnection() error {
	var finalErr error

	// if _, err := c.c.ReadFIFOQueue(1); err != nil {
	// 	finalErr = multierr.Append(finalErr, fmt.Errorf("ReadFIFOQueue error: %w", err))
	// }
	val, err := c.c.ReadDiscreteInputs(0, 10)
	glog.Infof("ReadDiscreteInputs(0, 10) = %v, %v", val, err)
	if err != nil {
		finalErr = multierr.Append(finalErr, fmt.Errorf("ReadDiscreteInputs error: %w", err))
	}
	time.Sleep(time.Millisecond * 300)

	const valuesPerRead = 1
	//for i := uint16(0); i < math.MaxUint16; i += 10 {
	for _, i := range []uint16{
		// Holding registers:
		28301,
		28302,
		28303,
		28306,
		28307,
		28308,
		28309,
		28310,
		28311,
		28312,
		28313,
		28314,
		28315,
		28316,
		28317,
		28318,
		28319,
		28320,
		28321,

		// Input registers:
		46801,
		46802,
		46803,
		46804,
		46805,
		46806,
		46807,
		46808,
		46809,
		46810,
	} {
		value, err := c.c.ReadHoldingRegisters(i, valuesPerRead)
		glog.Infof("ReadHoldingRegisters(%d, %d) = %v, %v", i, valuesPerRead, value, err)
		time.Sleep(time.Millisecond * 300)
		if err == nil {
			glog.Infof("SUCCCCCCCCCCCCCCCCCCCCCCCCCCCESSSSSSSSS ReadHoldingRegisters(%d, %d) = %v, %v", i, valuesPerRead, value, err)
		}
	}
	return fmt.Errorf("still testing")
	/*
			Not supported:

			if _, err := c.c.ReadCoils(1, 1); err != nil {
				finalErr = multierr.Append(finalErr, fmt.Errorf("ReadCoils error: %w", err))
			}
			if _, err := c.c.ReadFIFOQueue(1); err != nil {
				finalErr = multierr.Append(finalErr, fmt.Errorf("ReadFIFOQueue error: %w", err))
			}
			if _, err := c.c.ReadDiscreteInputs(1, 1); err != nil {
				finalErr = multierr.Append(finalErr, fmt.Errorf("ReadDiscreteInputs error: %w", err))
			}
		_, err := c.ReadState()
		if err != nil {
			finalErr = multierr.Append(finalErr, fmt.Errorf("ReadState failed: %w", err))
		}
	*/
	return finalErr
}

// State is a snapshot of the heat pump's state.
type State struct {
	collectionTime time.Time
	registerValues map[Register]uint16
}

// StateFromProto converts a state proto into a State object.
func StateFromProto(msg *chiltrix.State) (*State, error) {
	m := make(map[Register]uint16)
	for k, v := range msg.GetRegisterValues().GetHoldingRegisterValues() {
		if k > math.MaxUint16 {
			return nil, fmt.Errorf("register key %d is larger than the max uint16 %d", k, math.MaxInt16)
		}
		if v > math.MaxUint16 {
			return nil, fmt.Errorf("register[%d] value %d is larger than the max uint16 %d", k, v, math.MaxInt16)
		}
		m[Register(k)] = uint16(v)
	}
	if msg.GetCollectionTime() == nil {
		return nil, fmt.Errorf("state is missing collection time")
	}
	return &State{msg.CollectionTime.AsTime(), m}, nil
}

// Report returns a human readable summary of the state of the heat pump.
func (s *State) Report(omitZeros bool, interestingRegisters map[Register]bool) string {
	type entry struct {
		reg   Register
		value uint16
	}
	var entries []entry
	for k, v := range s.registerValues {
		if len(interestingRegisters) != 0 && !interestingRegisters[k] {
			continue
		}
		if (v == 0) && omitZeros {
			continue
		}
		entries = append(entries, entry{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].reg < entries[j].reg
	})
	b := &mdtable.Builder{}
	b.SetHeader([]string{"Register no.", "Register name", "Value"})
	for _, e := range entries {
		b.AddRow([]string{fmt.Sprintf("%d", e.reg), e.reg.String(), fmt.Sprintf("%d", e.value)})
	}
	return b.Build()
}

// String returns a human readable summary of the state of the heat pump.
func (s *State) String() string {
	return s.Report(false, nil)
}

// Proto returns the protobuf form of state.
func (s *State) Proto() *chiltrix.State {
	msg := &chiltrix.State{
		CollectionTime: timestamppb.New(s.collectionTime),
		RegisterValues: &chiltrix.RegisterSnapshot{
			HoldingRegisterValues: make(map[uint32]uint32),
		},
	}
	for reg, value := range s.registerValues {
		msg.GetRegisterValues().GetHoldingRegisterValues()[uint32(reg.uint16())] = uint32(value)
	}
	return msg
}

// CollectionTime returns the collection time of the heat pump state log entry.
func (s *State) CollectionTime() time.Time {
	return s.collectionTime
}

// RegisterValues returns a map of register values. The map should not be modified by the caller.
func (s *State) RegisterValues() map[Register]uint16 {
	return s.registerValues
}

// DiffStates returns a human readable description of the difference between two State values.
func DiffStates(a, b *State) (string, map[Register]bool) {
	diffRegs := map[Register]bool{}
	type entry struct {
		reg  Register
		a, b uint16
	}
	var entries []entry
	for k, v := range a.registerValues {
		if v != b.registerValues[k] {
			entries = append(entries, entry{k, v, b.registerValues[k]})
			diffRegs[k] = true
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].reg < entries[j].reg
	})
	builder := &mdtable.Builder{}
	builder.SetHeader([]string{"Register no.", "Register name", "Old Value", "New Value"})
	for _, e := range entries {
		builder.AddRow([]string{fmt.Sprintf("%d", e.reg), e.reg.String(), fmt.Sprintf("%d", e.a), fmt.Sprintf("%d", e.b)})
	}
	return builder.Build(), diffRegs
}

// Register is a modsbus register
type Register uint16

func (r Register) uint16() uint16 {
	return uint16(r)
}

// String returns a human-readable name of the modbus register.
func (r Register) String() string {
	return fmt.Sprintf("%d", r)
}
