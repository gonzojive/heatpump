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
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/goburrow/modbus"
	"github.com/goburrow/serial"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/mdtable"
	"github.com/gonzojive/heatpump/util/lockutil"
	"github.com/gonzojive/heatpump/util/modbusutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/gonzojive/heatpump/proto/fancoil"
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
	// If commands aren't spaced apart enough in time, the fan coil struggles.
	modbusCommandSpacing = time.Millisecond * 100
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
	c               modbus.Client
	modbusRTUCloser io.Closer
}

// Connect connects a new client to the heat pump or returns an error.
func Connect(ctx context.Context, p *Params) (*Client, error) {
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
	c := &Client{
		modbusutil.ClientWithLock(client, lockutil.WithGuaranteedTimeSinceLastRelease(&sync.Mutex{}, modbusCommandSpacing)),
		handler,
	}

	if err := c.CheckConnection(ctx); err != nil {
		handler.Close()
		return nil, err
	}
	return c, nil
}

// Close closes the modbus connection and frees up resources associated with the client.
func (c *Client) Close() error {
	return c.modbusRTUCloser.Close()
}

const (
	holdingRegisterStart = uint16(pb.RegisterName_REGISTER_NAME_ON_OFF)
	holdingRegisterCount = uint16(pb.RegisterName_REGISTER_NAME_UNIT_ADDRESS) - uint16(holdingRegisterStart) + 1

	inputRegisterStart = uint16(pb.RegisterName_REGISTER_NAME_ROOM_TEMPERATURE)
	inputRegisterCount = uint16(pb.RegisterName_REGISTER_NAME_COIL_TEMPERATURE_SENSOR_FAULT) - uint16(inputRegisterStart) + 1
)

// GetState returns a snapshot of the state of a single fan coil unit.
func (c *Client) GetState(_ context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	s, err := c.readRawState()
	if err != nil {
		return nil, err
	}
	return s.ResponseProto(), nil
}

func (c *Client) readRawState() (*State, error) {
	rawProto := &pb.RawRegisterSnapshot{
		RawValues: make(map[uint32]uint32, inputRegisterCount+holdingRegisterCount),
	}
	populateRawProto := func(registervalues []byte, startReg uint16) {
		for i := 0; i < len(registervalues)/2; i++ {
			value := uint16(registervalues[i*2])<<8 + uint16(registervalues[i*2+1])
			rawProto.RawValues[uint32(startReg)+uint32(i)] = uint32(value)
		}
	}

	holdingRegisterValues, err := c.c.ReadHoldingRegisters(holdingRegisterStart, holdingRegisterCount)
	if err != nil {
		return nil, grpc.Errorf(codes.Unavailable, "failed to read holding registers using modbus: %v", err)
	}
	populateRawProto(holdingRegisterValues, holdingRegisterStart)

	inputRegisterValues, err := c.c.ReadInputRegisters(inputRegisterStart, inputRegisterCount)
	if err != nil {
		return nil, grpc.Errorf(codes.Unavailable, "failed to read input registers using modbus: %v", err)
	}
	populateRawProto(inputRegisterValues, inputRegisterStart)

	return StateFromSnapshotProto(time.Now(), rawProto)
}

func (c *Client) setRegisterValue(reg Register, value uint16) error {
	// The fancoil doesn't support the WriteSingleRegister function.
	res, err := c.c.WriteMultipleRegisters(reg.uint16(), 1, uint16ToBytes(value))

	if err != nil {
		return fmt.Errorf("WriteSingleRegister error: %w (returned bytes %v)", err, res)
	}
	glog.Infof("set fancoil register value %d to %d; got response %v", reg, value, res)
	return nil
}

func uint16ToBytes(value uint16) []byte {
	out := make([]byte, 2)
	binary.BigEndian.PutUint16(out, value)
	return out
}

// CheckConnection attempts to connect to the heat pump and returns an error if the connection fails.
func (c *Client) CheckConnection(ctx context.Context) error {
	if _, err := c.GetState(ctx, &pb.GetStateRequest{
		FancoilName: fmt.Sprintf("%d", slaveID),
	}); err != nil {
		return fmt.Errorf("error getting state of fancoil; check connection: %w", err)
	}
	return nil
}

// State is a snapshot of the heat pump's state.
type State struct {
	collectionTime time.Time
	registerValues map[Register]uint16
	parsedState    *pb.State
}

// StateFromSnapshotProto converts a state proto into a State object.
func StateFromSnapshotProto(collectionTime time.Time, msg *pb.RawRegisterSnapshot) (*State, error) {
	m := make(map[Register]uint16)
	for k, v := range msg.GetRawValues() {
		if k > math.MaxUint16 {
			return nil, fmt.Errorf("register key %d is larger than the max uint16 %d", k, math.MaxInt16)
		}
		if v > math.MaxUint16 {
			return nil, fmt.Errorf("register[%d] value %d is larger than the max uint16 %d", k, v, math.MaxInt16)
		}
		m[Register(k)] = uint16(v)
	}
	parsed, err := parseRegisterValues(m)
	if err != nil {
		return nil, fmt.Errorf("error parsing register values: %v", err)
	}
	parsed.SnapshotTime = timestamppb.New(collectionTime)
	return &State{collectionTime, m, parsed}, nil
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
	parsedStr := ""
	if parsed, err := parseRegisterValues(s.registerValues); err != nil {
		parsedStr = fmt.Sprintf("error parsing register values: %v", err)
	} else {
		parsedStr = fmt.Sprintf("State proto: %s", prototext.Format(parsed))
	}
	return fmt.Sprintf("%s\n\n%s", b.Build(), parsedStr)
}

// String returns a human readable summary of the state of the heat pump.
func (s *State) String() string {
	return s.Report(false, nil)
}

// ResponseProto returns the protobuf form of state.
func (s *State) ResponseProto() *pb.GetStateResponse {
	msg := &pb.GetStateResponse{
		RawSnapshot: &pb.RawRegisterSnapshot{
			RawValues: make(map[uint32]uint32),
		},
		State: s.parsedState,
	}
	m := msg.GetRawSnapshot().GetRawValues()
	for reg, value := range s.registerValues {
		m[uint32(reg.uint16())] = uint32(value)
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
