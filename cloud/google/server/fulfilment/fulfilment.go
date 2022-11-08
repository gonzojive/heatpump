// Package fulfillment implements a server that responds to Google Smart
// Home "fulfilment" requests.
package fulfilment

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/proto/fancoil"
	smarthome "github.com/rmrobinson/google-smart-home-action-go"
)

const fixmeUserID = "4242"

type Service struct {
	sh *smarthome.Service
}

func NewService(psc *pubsub.Client) *Service {
	return &Service{
		sh: smarthome.NewService(zap.L(), &validator{}, &fulfilmentService{
			commandPublisher: newCommandPublisher(psc),
			lights: map[string]lightbulb{
				"abc": {
					id:         "abc",
					name:       "ABC Light",
					isOn:       true,
					brightness: 8,
					color: struct {
						hue        float64
						saturation float64
						value      float64
					}{
						hue:        .3,
						saturation: .6,
						value:      .3,
					},
				},
			},
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

type lightbulb struct {
	id         string
	name       string
	isOn       bool
	brightness int

	color struct {
		hue        float64
		saturation float64
		value      float64
	}
}

func (l *lightbulb) GetState() smarthome.DeviceState {
	return smarthome.NewDeviceState(true).
		RecordOnOff(l.isOn).
		RecordBrightness(l.brightness).
		RecordColorHSV(l.color.hue, l.color.saturation, l.color.value)
}

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
	return temperature(fDegrees-32) * 5 / 9
}

func (t temperature) Celcius() float64    { return float64(t) }
func (t temperature) Fahrenheit() float64 { return (t.Celcius() * 9 / 5) + 32 }

func (dev *fanCoilUnit) GetState() smarthome.DeviceState {
	// See https://developers.google.com/assistant/smarthome/guides/thermostat#sample-query-response
	return smarthome.NewDeviceState(true).
		RecordCustomState("thermostatMode", dev.mode).
		RecordCustomState("thermostatTemperatureSetpoint", dev.thermostatTemperatureSetpoint.Celcius()).
		RecordCustomState("thermostatTemperatureAmbient", dev.thermostatTemperatureAmbient.Celcius())
	//RecordCustomState("thermostatHumidityAmbient", dev.thermostatHumidityAmbient.Fahrenheit())
}

type fulfilmentService struct {
	commandPublisher *commandPublisher

	lights      map[string]lightbulb
	fanCoilUnit map[string]fanCoilUnit
}

func (srv *fulfilmentService) Sync(context.Context, string) (*smarthome.SyncResponse, error) {
	glog.Infof("sync")

	resp := &smarthome.SyncResponse{}
	for _, light := range srv.lights {
		ad := smarthome.NewLight(light.id)
		ad.Name = smarthome.DeviceName{
			DefaultNames: []string{
				"Test lamp",
			},
			Name: light.name,
		}
		ad.WillReportState = false
		ad.RoomHint = "test room"
		ad.DeviceInfo = smarthome.DeviceInfo{
			Manufacturer: "faltung systems",
			Model:        "tl001",
			HwVersion:    "0.2",
			SwVersion:    "0.3",
		}
		ad.
			AddOnOffTrait(false, false).
			AddBrightnessTrait(false).
			AddColourTrait(smarthome.HSV, false)

		resp.Devices = append(resp.Devices, ad)
	}

	for _, fcu := range srv.fanCoilUnit {
		// Fan coil units have built-in thermostats.
		dev := smarthome.NewDevice(fcu.id, "action.devices.types.THERMOSTAT")
		dev.AddCustomTrait("action.devices.traits.TemperatureSetting", func(setAttribute func(key string, val interface{})) {
			setAttribute("availableThermostatModes", []string{
				"off",
				"heat",
				"cool",
				"on",
			})
			setAttribute("thermostatTemperatureRange", map[string]int{
				"minThresholdCelsius": 15,
				"maxThresholdCelsius": 30,
			})
			setAttribute("thermostatTemperatureUnit", "F")
		})
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

	return resp, nil
}

func (es *fulfilmentService) Disconnect(context.Context, string) error {
	glog.Infof("disconnect")
	return nil
}

func (es *fulfilmentService) Query(_ context.Context, req *smarthome.QueryRequest) (*smarthome.QueryResponse, error) {
	glog.Infof("query")

	resp := &smarthome.QueryResponse{
		States: map[string]smarthome.DeviceState{},
	}

	for _, deviceArg := range req.Devices {
		if light, found := es.lights[deviceArg.ID]; found {
			resp.States[deviceArg.ID] = light.GetState()
		}
		if dev, found := es.fanCoilUnit[deviceArg.ID]; found {
			resp.States[deviceArg.ID] = dev.GetState()
		}
	}

	return resp, nil
}
func (es *fulfilmentService) Execute(ctx context.Context, req *smarthome.ExecuteRequest) (*smarthome.ExecuteResponse, error) {
	glog.Infof("Execute")
	resp := &smarthome.ExecuteResponse{
		UpdatedState: smarthome.NewDeviceState(true),
	}
	for _, cmd := range req.Commands {
		for _, target := range cmd.TargetDevices {
			fcu, ok := es.fanCoilUnit[target.ID]
			if !ok {
				continue
			}
			for _, cmdObj := range cmd.Commands {
				parsedCmd, err := parseTemperatureSettingCommand(cmdObj.Generic)
				if err != nil {
					return nil, fmt.Errorf("failed to parse command for %q: %w", target.ID, err)
				}

				switch c := parsedCmd.(type) {
				// https://developers.google.com/assistant/smarthome/traits/temperaturesetting#action.devices.commands.thermostattemperaturesetpoint
				case setSetpointCommand:
					fcu.thermostatTemperatureSetpoint = c.Setpoint
					glog.Infof("setpoint adjusted to %f", c.Setpoint.Fahrenheit())
					resp.UpdatedDevices = append(resp.UpdatedDevices, fcu.id)
					es.commandPublisher.executeFanCoilCommand(ctx, &fancoil.SetStateRequest{
						HeatingTargetTemperature: &fancoil.Temperature{
							DegreesCelcius: float32(c.Setpoint.Celcius()),
						},
					})
				}
			}
		}
	}
	return resp, nil
}

func parseTemperatureSettingCommand(cmdObj *smarthome.CommandGeneric) (any, error) {
	switch cmdObj.Command {
	// https://developers.google.com/assistant/smarthome/traits/temperaturesetting#action.devices.commands.thermostattemperaturesetpoint
	case "action.devices.commands.ThermostatTemperatureSetpoint":
		val := cmdObj.Params["thermostatTemperatureSetpoint"]
		switch val := val.(type) {
		case float64:
			return setSetpointCommand{
				Setpoint: temperature(val),
			}, nil
		default:
			return nil, fmt.Errorf("invalid request has bad thermostatTemperatureSetpoint param %v", val)
		}
	default:
		return nil, fmt.Errorf("couldn't parse command %q: %+v", cmdObj.Command, cmdObj.Params)
	}
}

type setSetpointCommand struct {
	// Target temperature setpoint. Supports up to one decimal place.
	Setpoint temperature `json:"thermostatTemperatureSetpoint`
}
