// This file contains struct definitions for Googlde Device Traits related to
// the fan coil units sold by Chiltrix.

package fulfilment

import (
	"github.com/gonzojive/heatpump/proto/fancoil"
	smarthome "github.com/rmrobinson/google-smart-home-action-go"
)

func fanSettingFromName(name string) fancoil.FanSetting {
	return fancoil.FanSetting(fancoil.FanSetting_value[name])
}

func fanSettingToName(fs fancoil.FanSetting) string {
	return fs.String()
}

var fancoilFanSpeedAttribute = FanSpeedAttributes{
	Reversible:              false,
	CommandOnlyFanSpeed:     false,
	SupportsFanSpeedPercent: false,
	AvailableFanSpeeds: AvailableFanSpeeds{
		Ordered: true,
		Speeds: []FanSpeedSetting{
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_OFF),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"Off"},
						Lang:         "en",
					},
				},
			},
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_ULTRA_LOW),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"Ultra-low"},
						Lang:         "en",
					},
				},
			},
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_LOW),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"Low"},
						Lang:         "en",
					},
				},
			},
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_MEDIUM),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"Medium"},
						Lang:         "en",
					},
				},
			},
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_HIGH),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"High"},
						Lang:         "en",
					},
				},
			},
			{
				SpeedName: fanSettingToName(fancoil.FanSetting_FAN_SETTING_AUTO),
				SpeedValues: []FanSpeedValue{
					{
						SpeedSynonym: []string{"auto"},
						Lang:         "en",
					},
				},
			},
		},
	},
}

// ID of a device trait name, like action.devices.traits.FanSpeed.
type TraitName string

// Attributes for the FanSpeed attribute, "action.devices.traits.FanSpeed".
//
// Described here: https://developers.google.com/assistant/smarthome/traits/fanspeed.
type FanSpeedAttributes struct {
	// If set to true, this device supports blowing the fan in both directions
	// and can accept the command to reverse fan direction.
	Reversible bool `json:"reversible"`

	// Indicates if the device supports using one-way (true) or two-way (false)
	// communication. Set this attribute to true if the device cannot respond to
	// a QUERY intent or Report State for this trait.
	CommandOnlyFanSpeed bool `json:"commandOnlyFanSpeed"`

	// If set to true, this device will accept commands for adjusting the speed
	// using a percentage from 0.0 to 100.0.
	SupportsFanSpeedPercent bool `json:"supportsFanSpeedPercent"`

	// List of available fan speeds.
	AvailableFanSpeeds AvailableFanSpeeds `json:"availableFanSpeeds"`
}

// AddToDevice adds the attributes to the given Device object.
func (a *FanSpeedAttributes) AddToDevice(dev *smarthome.Device) {
	dev.Traits[string(a.TraitName())] = true

	dev.Attributes["reversible"] = a.Reversible
	dev.Attributes["commandOnlyFanSpeed"] = a.CommandOnlyFanSpeed
	dev.Attributes["supportsFanSpeedPercent"] = a.SupportsFanSpeedPercent
	dev.Attributes["availableFanSpeeds"] = a.AvailableFanSpeeds
}

func (a *FanSpeedAttributes) TraitName() TraitName {
	return "action.devices.traits.FanSpeed"
}

type AvailableFanSpeeds struct {
	Speeds  []FanSpeedSetting `json:"speeds"`
	Ordered bool              `json:"ordered"`
}

type FanSpeedSetting struct {
	SpeedName   string          `json:"speed_name"`
	SpeedValues []FanSpeedValue `json:"speed_values"`
}

type FanSpeedValue struct {
	SpeedSynonym []string `json:"speed_synonym"`
	Lang         string   `json:"lang"`
}

type FanSpeedState struct {
	CurrentFanSpeedSetting string `json:"currentFanSpeedSetting"`
	CurrentFanSpeedPercent int    `json:"currentFanSpeedPercent"`
}

type SetFanSpeedNameCommand struct {
	// The request fan speed name.
	FanSpeed string `json:"fanSpeed"`
}

type SetFanSpeedPercentCommand struct {
	// The requested speed setting percentage.
	FanSpeedPercent float32 `json:"fanSpeedPercent"`
}
