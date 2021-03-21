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

func parseRegisterValues(values map[Register]uint16) (*pb.State, error) {
	out := &pb.State{
		ModbusAddress: uint32(values[Register(pb.RegisterName_REGISTER_NAME_UNIT_ADDRESS)]),
	}

	// Temperature values
	for _, spec := range []struct {
		reg Register
		dst **pb.Temperature
	}{
		{
			Register(pb.RegisterName_REGISTER_NAME_ROOM_TEMPERATURE),
			&out.RoomTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_COIL_TEMPERATURE),
			&out.CoilTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_HEATING_SET_TEMPERATURE),
			&out.HeatingTargetTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_COOLING_SET_TEMPERATURE),
			&out.CoolingTargetTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_ANTI_COOLING_WIND_SETTING_TEMPERATURE),
			&out.AntiCoolingTargetTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_COOLING_SET_TEMPERATURE_AUTO),
			&out.AutoModeCoolingTargetTemperature,
		},
		{
			Register(pb.RegisterName_REGISTER_NAME_HEATING_SET_TEMPERATURE_AUTO),
			&out.AutoModeHeatingTargetTemperature,
		},
	} {
		if val, ok := values[spec.reg]; ok {
			t, err := parseTemp(val)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s value %v: %w", spec.reg, val, err)
			}
			*spec.dst = t
		}
	}

	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_CURRENT_FAN_SPEED)]; ok {
		fs, err := parseFanSetting(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse current fan speed value %v: %w", val, err)
		}
		out.CurrentFanSetting = fs
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_FANSPEED)]; ok {
		fs, err := parseFanSetting(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse current fan speed value %v: %w", val, err)
		}
		out.PreferenceFanSetting = fs
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_USE_FLOOR_HEATING)]; ok {
		enumVal, err := parseFloorHeatingMode(val)
		if err != nil {
			return nil, err
		}
		out.FloorHeatingMode = enumVal
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_USE_FAHRENHEIT)]; ok {
		enumVal, err := parseTempDisplayUnit(val)
		if err != nil {
			return nil, err
		}
		out.DisplayTemperatureUnits = enumVal
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_VALVE_ON_OFF)]; ok {
		enumVal, err := parseValveState(val)
		if err != nil {
			return nil, err
		}
		out.ValveState = enumVal
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_FAN_RPM)]; ok {
		out.FanSpeed = &pb.FanSpeed{Rpm: int64(val)}
	}

	return out, nil
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

func parseFanSetting(value uint16) (pb.FanSetting, error) {
	switch value {
	case 0:
		return pb.FanSetting_FAN_SETTING_AUTO, nil
	default:
		if _, known := pb.FanSetting_name[int32(value)]; known {
			return pb.FanSetting(value), nil
		}
		return pb.FanSetting(value), fmt.Errorf("unknown fan setting value %v", value)
	}
}

func parseTempDisplayUnit(value uint16) (pb.TemperatureUnits, error) {
	switch value {
	case 0:
		return pb.TemperatureUnits_TEMPERATURE_UNITS_CELCIUS, nil
	case 1:
		return pb.TemperatureUnits_TEMPERATURE_UNITS_FAHRENHEIT, nil
	default:
		return pb.TemperatureUnits_TEMPERATURE_UNITS_UNSPECIFIED, fmt.Errorf("invalid temperature unit value %d", value)
	}
}

func parseFloorHeatingMode(value uint16) (pb.FloorHeatingMode, error) {
	switch value {
	case 0:
		return pb.FloorHeatingMode_FLOOR_HEATING_MODE_OFF, nil
	case 1:
		return pb.FloorHeatingMode_FLOOR_HEATING_MODE_ON, nil
	default:
		return pb.FloorHeatingMode_FLOOR_HEATING_MODE_UNSPECIFIED, fmt.Errorf("invalid floor heating mode value %d", value)
	}
}

func parseValveState(value uint16) (pb.ValveState, error) {
	switch value {
	case 0:
		return pb.ValveState_VALVE_STATE_OFF, nil
	case 1:
		return pb.ValveState_VALVE_STATE_ON, nil
	default:
		return pb.ValveState_VALVE_STATE_UNSPECIFIED, fmt.Errorf("invalid floor heating mode value %d", value)
	}
}
