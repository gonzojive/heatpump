package fancoil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/gonzojive/heatpump/proto/fancoil"
)

func Test_parseRegisterValues(t *testing.T) {
	tests := []struct {
		name      string
		rawValues map[Register]uint16
		want      *pb.State
		wantErr   bool
	}{
		{
			"real example 1",
			map[Register]uint16{
				Register(pb.RegisterName_REGISTER_NAME_ON_OFF):   1,
				Register(pb.RegisterName_REGISTER_NAME_MODE):     4,
				Register(pb.RegisterName_REGISTER_NAME_FANSPEED): 2,
				Register(28304): 0,
				Register(28305): 0,
				Register(pb.RegisterName_REGISTER_NAME_TIMER_OFF1):                            0,
				Register(pb.RegisterName_REGISTER_NAME_TIMER_OFF2):                            0,
				Register(pb.RegisterName_REGISTER_NAME_MAX_SET_TEMPERATURE):                   300,
				Register(pb.RegisterName_REGISTER_NAME_MIN_SET_TEMPERATURE):                   80,
				Register(pb.RegisterName_REGISTER_NAME_COOLING_SET_TEMPERATURE):               210,
				Register(pb.RegisterName_REGISTER_NAME_HEATING_SET_TEMPERATURE):               190,
				Register(pb.RegisterName_REGISTER_NAME_COOLING_SET_TEMPERATURE_AUTO):          260,
				Register(pb.RegisterName_REGISTER_NAME_HEATING_SET_TEMPERATURE_AUTO):          200,
				Register(pb.RegisterName_REGISTER_NAME_ANTI_COOLING_WIND_SETTING_TEMPERATURE): 250,
				Register(pb.RegisterName_REGISTER_NAME_START_ANTI_HOT_WIND):                   0,
				Register(pb.RegisterName_REGISTER_NAME_START_ULTRA_LOW_WIND):                  0,
				Register(pb.RegisterName_REGISTER_NAME_USE_VALVE):                             0,
				Register(pb.RegisterName_REGISTER_NAME_USE_FLOOR_HEATING):                     0,
				Register(pb.RegisterName_REGISTER_NAME_USE_FAHRENHEIT):                        1,
				Register(pb.RegisterName_REGISTER_NAME_MASTER_SLAVE):                          1,
				Register(pb.RegisterName_REGISTER_NAME_UNIT_ADDRESS):                          15,
				Register(pb.RegisterName_REGISTER_NAME_ROOM_TEMPERATURE):                      200,
				Register(pb.RegisterName_REGISTER_NAME_COIL_TEMPERATURE):                      410,
				Register(pb.RegisterName_REGISTER_NAME_CURRENT_FAN_SPEED):                     2,
				Register(pb.RegisterName_REGISTER_NAME_FAN_RPM):                               1005,
				Register(pb.RegisterName_REGISTER_NAME_VALVE_ON_OFF):                          0,
				Register(pb.RegisterName_REGISTER_NAME_REMOTE_ON_OFF):                         0,
				Register(pb.RegisterName_REGISTER_NAME_SIMULATION_SIGNAL):                     1,
				Register(pb.RegisterName_REGISTER_NAME_FAN_SPEED_SIGNAL_FEEDBACK_FAULT):       0,
				Register(pb.RegisterName_REGISTER_NAME_ROOM_TEMPERATURE_SENSOR_FAULT):         0,
				Register(pb.RegisterName_REGISTER_NAME_COIL_TEMPERATURE_SENSOR_FAULT):         1,
			},
			&pb.State{
				HeatingTargetTemperature:         &pb.Temperature{DegreesCelcius: 19},
				CoolingTargetTemperature:         &pb.Temperature{DegreesCelcius: 21},
				AntiCoolingTargetTemperature:     &pb.Temperature{DegreesCelcius: 25},
				AutoModeHeatingTargetTemperature: &pb.Temperature{DegreesCelcius: 20},
				AutoModeCoolingTargetTemperature: &pb.Temperature{DegreesCelcius: 26},

				FloorHeatingMode:        pb.FloorHeatingMode_FLOOR_HEATING_MODE_OFF,
				DisplayTemperatureUnits: pb.TemperatureUnits_TEMPERATURE_UNITS_FAHRENHEIT,

				CoilTemperature:      &pb.Temperature{DegreesCelcius: 41},
				RoomTemperature:      &pb.Temperature{DegreesCelcius: 20},
				PreferenceFanSetting: pb.FanSetting_FAN_SETTING_LOW,
				CurrentFanSetting:    pb.FanSetting_FAN_SETTING_LOW,
				FanSpeed:             &pb.FanSpeed{Rpm: 1005},
				ModbusAddress:        &pb.ModbusAddress{Address: 15},
				ValveState:           pb.ValveState_VALVE_STATE_OFF,
				ValveSetting:         pb.ValveSetting_VALVE_SETTING_DISABLED,
				PowerStatus:          pb.PowerStatus_POWER_STATUS_ON,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRegisterValues(tt.rawValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRegisterValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("parseRegisterValues got unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_parseTemperature(t *testing.T) {
	tests := []struct {
		name          string
		registerValue uint16
		want          *pb.Temperature
		wantErr       bool
	}{
		{
			name:          "30 degrees C",
			registerValue: 300,
			want:          &pb.Temperature{DegreesCelcius: 30},
		},
		{
			name:          "-30 degrees C",
			registerValue: 300 + negativeTempBit,
			want:          &pb.Temperature{DegreesCelcius: -30},
		},
		{
			name:          "0 degrees C",
			registerValue: 0,
			want:          &pb.Temperature{DegreesCelcius: 0},
		},
		{
			name:          "22.2 degrees C",
			registerValue: 222,
			want:          &pb.Temperature{DegreesCelcius: 22.2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTemp(tt.registerValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRegisterValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("parseRegisterValues got unexpected diff (-want, +got):\n%s", diff)
			}
			gotRegister, err := tempToModbusValue(tt.want)
			if err != nil {
				t.Errorf("valid temperature got error converting to modbus value %v", err)
			}
			if gotRegister != tt.registerValue {
				t.Errorf("temperature %v converted to %d, expected %d", tt.want, gotRegister, tt.registerValue)
			}
		})
	}
}
