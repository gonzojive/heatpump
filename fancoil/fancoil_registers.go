package fancoil

import (
	"fmt"
	"math"

	pb "github.com/gonzojive/heatpump/proto/fancoil"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	faultyTemperatureCode = 32767
	negativeTempBit       = uint16(1 << 15)
	tempNonSignMask       = 0xFFFF ^ negativeTempBit
)

// Enum maps
var (
	valveSettingMap = parseEnumRegisterValues(pb.ValveSetting(0).Descriptor())
	fanSettingDef   = parseEnumRegisterValues(pb.FanSetting(0).Descriptor())
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
	out := &pb.State{}

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

	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_UNIT_ADDRESS)]; ok {
		out.ModbusAddress = &pb.ModbusAddress{Address: uint32(val)}
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
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_ON_OFF)]; ok {
		enumVal, err := parsePowerStatus(val)
		if err != nil {
			return nil, err
		}
		out.PowerStatus = enumVal
	}
	if val, ok := values[Register(pb.RegisterName_REGISTER_NAME_FAN_RPM)]; ok {
		out.FanSpeed = &pb.FanSpeed{Rpm: int64(val)}
	}
	valveSetting, err := valveSettingMap.parseRegister(values, Register(pb.RegisterName_REGISTER_NAME_USE_VALVE))
	if err != nil {
		return nil, err
	}
	out.ValveSetting = pb.ValveSetting(valveSetting)

	FanSetting, err := fanSettingDef.parseRegister(values, Register(pb.RegisterName_REGISTER_NAME_FANSPEED))
	if err != nil {
		return nil, err
	}
	out.PreferenceFanSetting = pb.FanSetting(FanSetting)

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

func fanSettingToModbusValue(v pb.FanSetting) (uint16, error) {
	switch v {
	case pb.FanSetting_FAN_SETTING_UNSPECIFIED:
		return 0, fmt.Errorf("no modbus value for %s", v)
	case pb.FanSetting_FAN_SETTING_AUTO:
		return 0, nil
	default:
		return uint16(v), nil
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

func parsePowerStatus(value uint16) (pb.PowerStatus, error) {
	switch value {
	case 0:
		return pb.PowerStatus_POWER_STATUS_OFF, nil
	case 1:
		return pb.PowerStatus_POWER_STATUS_ON, nil
	default:
		return pb.PowerStatus_POWER_STATUS_UNSPECIFIED, fmt.Errorf("invalid floor heating mode value %d", value)
	}
}

func powerStatusToModbusValue(v pb.PowerStatus) (uint16, error) {
	switch v {
	case pb.PowerStatus_POWER_STATUS_UNSPECIFIED:
		return 0, fmt.Errorf("no modbus value for %s", v)
	case pb.PowerStatus_POWER_STATUS_OFF:
		return 0, nil
	case pb.PowerStatus_POWER_STATUS_ON:
		return 1, nil
	default:
		return 0, fmt.Errorf("unsupported power status value %s", v)
	}
}

type enumRegisterValueMap struct {
	numberToValue map[protoreflect.EnumNumber]uint16
}

func (m *enumRegisterValueMap) getByRegisterValue(v uint16) (protoreflect.EnumNumber, bool) {
	for number, value := range m.numberToValue {
		if value == v {
			return number, true
		}
	}
	return 0, false
}

func (m *enumRegisterValueMap) getByEnumNumber(n protoreflect.EnumNumber) (uint16, bool) {
	v, ok := m.numberToValue[n]
	return v, ok
}

func (m *enumRegisterValueMap) parseRegister(registerValues map[Register]uint16, register Register) (protoreflect.EnumNumber, error) {
	v, registerIsSet := registerValues[register]
	if !registerIsSet {
		// Zero is always UNSPECIFIED.
		return 0, nil
	}
	if num, ok := m.getByRegisterValue(v); ok {
		return num, nil
	}
	return 0, fmt.Errorf("unknown register value %d does not correspond to any enum value in the map: %v", v, m.numberToValue)
}

func parseEnumRegisterValues(def protoreflect.EnumDescriptor) *enumRegisterValueMap {
	m := make(map[protoreflect.EnumNumber]uint16)
	for i := 0; i < def.Values().Len(); i++ {
		valDef := def.Values().Get(i)
		if !proto.HasExtension(valDef.Options(), pb.E_ModbusOptions) {
			continue
		}
		mbOpt := proto.GetExtension(valDef.Options(), pb.E_ModbusOptions).(*pb.ModbusEnumValueOptions)
		if mbOpt.GetRegisterValue() < 0 || mbOpt.GetRegisterValue() > math.MaxUint16 {
			panic(fmt.Errorf("invalid register value for field %v: not a uint16", valDef))
		}
		m[valDef.Number()] = uint16(mbOpt.GetRegisterValue())
	}
	return &enumRegisterValueMap{m}
}
