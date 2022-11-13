// Package fulfillment implements a server that responds to Google Smart
// Home "fulfilment" requests.
package fulfilment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/proto/fancoil"
	"github.com/gonzojive/heatpump/util/must"
	smarthome "github.com/rmrobinson/google-smart-home-action-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/prototext"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

const fixmeUserID = "4242"

type Service struct {
	sh *smarthome.Service
}

func NewService(psc *pubsub.Client, ssc cpb.StateServiceClient) *Service {
	return &Service{
		sh: smarthome.NewService(zap.L(), &validator{}, &fulfilmentService{
			ssc:              ssc,
			commandPublisher: newCommandPublisher(psc),
			fanCoilUnit: map[string]fanCoilUnit{
				"xyz": {
					id:                            "xyz",
					name:                          "XYZ Fan Coil Unit",
					mode:                          "heat",
					thermostatTemperatureSetpoint: tempFromFahrenheit(72),
					thermostatTemperatureAmbient:  tempFromFahrenheit(66),
				},
			},
		}, nil),
	}
}

// GoogleFulfillmentHandler returns a handler for fulfillment requests.
func (s *Service) GoogleFulfillmentHandler() http.Handler {
	return http.HandlerFunc(s.sh.GoogleFulfillmentHandler)
}

type validator struct{}

// Validate performs the actual token validation. Returning an error will force
// validation to fail.
//
// The user ID that corresponds to the token should be returned on success.
func (v *validator) Validate(ctx context.Context, token string) (string, error) {
	return fixmeUserID, nil
}

var _ smarthome.AccessTokenValidator = (*validator)(nil)

type fanCoilUnit struct {
	id                            string
	name                          string
	mode                          string
	thermostatTemperatureSetpoint temperature
	thermostatTemperatureAmbient  temperature
	thermostatHumidityAmbient     temperature
}

type temperature float64

func tempFromFahrenheit(fDegrees float64) temperature {
	return tempFromCelcius(fDegrees-32) * 5 / 9
}

func tempFromCelcius(fDegrees float64) temperature {
	return temperature(fDegrees)
}

func (t temperature) Celcius() float64    { return float64(t) }
func (t temperature) Fahrenheit() float64 { return (t.Celcius() * 9 / 5) + 32 }

func (dev *fanCoilUnit) GetState(ctx context.Context, ssc cpb.StateServiceClient) (*smarthome.DeviceState, error) {

	stateResp, err := ssc.GetDeviceState(ctx, &cpb.GetDeviceStateRequest{
		Name: "", //dev.id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve state of device %q from StateService: %w", dev.id, err)
	}

	// See https://developers.google.com/assistant/smarthome/guides/thermostat#sample-query-response
	s := smarthome.NewDeviceState(true)

	switch stateResp.GetFancoilState().GetMode() {
	case fancoil.Mode_MODE_COOLING:
		s.State["thermostatMode"] = "cool"
		s.State["thermostatTemperatureSetpoint"] = stateResp.GetFancoilState().GetCoolingTargetTemperature().GetDegreesCelcius()
	case fancoil.Mode_MODE_HEATING:
		s.State["thermostatMode"] = "heat"
		s.State["thermostatTemperatureSetpoint"] = stateResp.GetFancoilState().GetHeatingTargetTemperature().GetDegreesCelcius()

	default:
		s.State["thermostatMode"] = "heat"
		s.State["thermostatTemperatureSetpoint"] = stateResp.GetFancoilState().GetHeatingTargetTemperature().GetDegreesCelcius()
	}

	switch setting := stateResp.GetFancoilState().GetCurrentFanSetting(); setting {
	case fancoil.FanSetting_FAN_SETTING_UNSPECIFIED:
		glog.Errorf("unknown fan speed state %v", stateResp.GetFancoilState().GetCurrentFanSetting())
	default:
		s.State["currentFanSpeedSetting"] = fanSettingToName(setting)
	}

	s.State["thermostatTemperatureAmbient"] = stateResp.GetFancoilState().GetRoomTemperature().GetDegreesCelcius()

	glog.Infof("state of fan coil unit: %s", prototext.Format(stateResp))

	return &s, nil
}

type fulfilmentService struct {
	commandPublisher *commandPublisher

	fanCoilUnit map[string]fanCoilUnit

	ssc cpb.StateServiceClient
}

// Returns the list of devices associated with the given user and their capabilities.
func (fs *fulfilmentService) Sync(ctx context.Context, userID string) (*smarthome.SyncResponse, error) {
	glog.Infof("sync")

	resp := &smarthome.SyncResponse{}

	glog.Infof("got request to sync for user id %q", userID)

	for _, fcu := range fs.fanCoilUnit {
		// Fan coil units have built-in thermostats.
		//dev := smarthome.NewDevice(fcu.id, "action.devices.types.THERMOSTAT")
		dev := smarthome.NewDevice(fcu.id, "action.devices.types.THERMOSTAT")

		dev.Traits["action.devices.traits.TemperatureSetting"] = true
		dev.Attributes["availableThermostatModes"] = []string{
			"off",
			"heat",
			"cool",
		}
		dev.Attributes["thermostatTemperatureRange"] = map[string]int{
			"minThresholdCelsius": 15,
			"maxThresholdCelsius": 30,
		}
		dev.Attributes["thermostatTemperatureUnit"] = "F"

		fancoilFanSpeedAttribute.AddToDevice(dev)

		dev.Name = smarthome.DeviceName{
			DefaultNames: []string{"Chiltrix Fan Coil Unit"},
			Name:         fcu.name,
			Nicknames:    []string{},
		}
		dev.DeviceInfo = smarthome.DeviceInfo{
			Manufacturer: "Chiltrix",
			Model:        "FCU1",
			HwVersion:    "0.2",
			SwVersion:    "0.3",
		}
		dev.WillReportState = true

		resp.Devices = append(resp.Devices, dev)
	}

	glog.Infof("sync response: %s", string(must.Value(json.MarshalIndent(resp, "", "  "))))

	return resp, nil
}

func (fs *fulfilmentService) Disconnect(context.Context, string) error {
	glog.Infof("disconnect")
	return nil
}

func (fs *fulfilmentService) Query(ctx context.Context, req *smarthome.QueryRequest) (*smarthome.QueryResponse, error) {
	glog.Infof("query")

	resp := &smarthome.QueryResponse{
		States: map[string]smarthome.DeviceState{},
	}

	for _, deviceArg := range req.Devices {

		glog.Infof("query for device %q...", deviceArg.ID)

		deviceName := deviceArg.ID
		if dev, found := fs.fanCoilUnit[deviceName]; found {
			state, err := dev.GetState(ctx, fs.ssc)
			if err != nil {
				glog.Warningf("query for device %q returned error: %v", deviceArg.ID, err)
				return nil, err
			} else {
				glog.Infof("got state for device %q: %s", deviceArg.ID, string(must.Value(json.MarshalIndent(*state, "", "  "))))
				resp.States[deviceName] = *state
			}
		} else {
			glog.Warningf("query for device %q didn't find device")
		}
	}

	return resp, nil
}
func (fs *fulfilmentService) Execute(ctx context.Context, req *smarthome.ExecuteRequest) (*smarthome.ExecuteResponse, error) {
	glog.Infof("Execute command %s", string(must.Value(json.MarshalIndent(req, "", "  "))))
	resp := &smarthome.ExecuteResponse{
		UpdatedState: smarthome.NewDeviceState(true),
	}
	for _, cmd := range req.Commands {
		for _, target := range cmd.TargetDevices {
			fcu, ok := fs.fanCoilUnit[target.ID]
			if !ok {
				continue
			}
			for _, cmdObj := range cmd.Commands {
				if parsed := cmdObj.SetFanSpeed; parsed != nil && parsed.FanSpeed != nil {
					setting := fanSettingFromName(*parsed.FanSpeed)
					glog.Infof("setting fan speed to to %q (%s)", *parsed.FanSpeed, setting)

					resp.UpdatedDevices = append(resp.UpdatedDevices, fcu.id)
					fs.commandPublisher.executeFanCoilCommand(ctx, &fancoil.SetStateRequest{
						PreferenceFanSetting: setting,
					})
				}
				if parsed := cmdObj.ThermostatTemperatureSetpoint; parsed != nil {
					// https://developers.google.com/assistant/smarthome/traits/temperaturesetting#action.devices.commands.thermostattemperaturesetpoint
					fcu.thermostatTemperatureSetpoint = tempFromCelcius(float64(parsed.ThermostatTemperatureSetpointCelcius))
					glog.Infof("setpoint adjusted to %f", fcu.thermostatTemperatureSetpoint.Fahrenheit())
					resp.UpdatedDevices = append(resp.UpdatedDevices, fcu.id)
					fs.commandPublisher.executeFanCoilCommand(ctx, &fancoil.SetStateRequest{
						HeatingTargetTemperature: &fancoil.Temperature{
							DegreesCelcius: float32(fcu.thermostatTemperatureSetpoint.Celcius()),
						},
					})
				}
			}
		}
	}
	return resp, nil
}

// func parseGenericCommand(cmdObj *smarthome.CommandGeneric) (any, error) {
// 	switch cmdObj.Command {
// 	// https://developers.google.com/assistant/smarthome/traits/temperaturesetting#action.devices.commands.thermostattemperaturesetpoint
// 	case "action.devices.commands.ThermostatTemperatureSetpoint":
// 		val := cmdObj.Params["thermostatTemperatureSetpoint"]
// 		switch val := val.(type) {
// 		case float64:
// 			return setSetpointCommand{
// 				Setpoint: temperature(val),
// 			}, nil
// 		default:
// 			return nil, fmt.Errorf("invalid request has bad thermostatTemperatureSetpoint param %v", val)
// 		}
// 	case "action.devices.commands.SetFanSpeed":
// 		jsonBytes, err := json.Marshal(cmdObj.Params)
// 		if err != nil {
// 			return nil, err
// 		}
// 		parsed := &SetFanSpeedNameCommand{}
// 		json.Unmarshal(data []byte, v any)
// 	default:
// 		return nil, fmt.Errorf("couldn't parse command %q: %+v", cmdObj.Command, cmdObj.Params)
// 	}
// }

// type setSetpointCommand struct {
// 	// Target temperature setpoint. Supports up to one decimal place.
// 	Setpoint temperature `json:"thermostatTemperatureSetpoint`
// }
