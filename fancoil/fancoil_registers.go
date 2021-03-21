package fancoil

import (
	"fmt"

	pb "github.com/gonzojive/heatpump/proto/fancoil"
)

const (
	faultyTemperatureCode = 32767
	negativeTempBit       = uint16(1 << 15)
	tempNonSignMask       = 0xFFFF ^ negativeTempBit
)

// Register is a modsbus register
type Register pb.RegisterName

func (r Register) uint16() uint16 {
	return uint16(r)
}

// String returns a human-readable name of the modbus register.
func (r Register) String() string {
	return pb.RegisterName(r).String()
}

func parseTemp(value uint16) (*pb.Temperature, error) {
	// Signed byte ，Precision 0.1℃，Formula：T*10，Temperature range ：-30~97℃
	// （if temperature is shown as 25°C, data transmission is 250 according to
	// the preceding formula. When bit15=1 , it means minus. when bit15=0, it
	// means integer );When this value is 32767, corresponding sensor is
	// faulty.）
	if value == faultyTemperatureCode {
		return nil, fmt.Errorf("temperature sensor fault, got code %d", value)
	}
	signedValue := int16(tempNonSignMask & value)
	if (value & negativeTempBit) != 0 {
		signedValue *= -1
	}

	return &pb.Temperature{DegreesCelcius: float32(signedValue) / 10}, nil
}

func parseRegisterValues(values map[Register]uint16) (*pb.State, error) {
	out := &pb.State{}

	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_ROOM_TEMPERATURE)]; ok {
		t, err := parseTemp(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse room temperature value %v: %w", val, err)
		}
		out.RoomTemperature = t
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_COIL_TEMPERATURE)]; ok {
		t, err := parseTemp(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coil temperature value %v: %w", val, err)
		}
		out.CoilTemperature = t
	}

	return out, nil
}
