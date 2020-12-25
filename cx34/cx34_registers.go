package cx34

import (
	"fmt"
	"time"

	"github.com/gonzojive/heatpump/units"
)

const (
	// Specific heat of the system depends on the glycol mix and temperature... but
	// for now we just use a constant.
	constantSpecificHeat units.SpecificHeat = units.KilojoulePerKilogramKelvin * 4.0
)

var (
	waterDensity = units.DensityFromRatio(units.Kilogram*997, units.CubicMeter)
)

// This file is used to assign names to modbus registers.

// FlowRate returns the water flow rate measured by the CX34's flow sensor.
//
// The flow sensor is made by the same company that makes this one:
// https://www.adafruit.com/product/828?gclid=Cj0KCQiAlZH_BRCgARIsAAZHSBmfM9AVkdnye4p7RVf_cbKDm6n6jILBT9ILjkvpg8PnLjz_38tU324aAsk0EALw_wcB
func (s *State) FlowRate() units.FlowRate {
	decilitersPerMinute := s.registerValues[WaterFlowRate]
	return units.LiterPerMinute.Scale(float64(decilitersPerMinute) / 10.0)
}

// SuctionTemp returns the "suction temperature" of the unit.
//
// A typical value in heating mode is 3.3 degrees C, so I'm not sure what this
// measures exactly.
func (s *State) SuctionTemp() units.Temperature {
	deciDegreesC := s.registerValues[SuctionTemp]
	return units.FromCelsius(float64(deciDegreesC) / 10.0)
}

// ACHeatingTargetTemp returns the active setpoint temperature.
//
// Currently this only returns the target AC temperature in heating mode. It
// should be updated to return the cooling mode temperature when in cooling
// mode.
func (s *State) ACHeatingTargetTemp() units.Temperature {
	degreesC := s.registerValues[TargetACHeatingModeTemp]
	return units.FromCelsius(float64(degreesC))
}

// ACOutletWaterTemp returns the temperature at the water outlet; values are -30~97℃.
func (s *State) ACOutletWaterTemp() units.Temperature {
	deciDegreesC := s.registerValues[ACOutletWaterTemp]
	return units.FromCelsius(float64(deciDegreesC) / 10.0)
}

// ACInletWaterTemp returns the temperature at the water inlet; values are -30~97℃.
func (s *State) ACInletWaterTemp() units.Temperature {
	deciDegreesC := s.registerValues[WaterInletSensorTemp1]
	return units.FromCelsius(float64(deciDegreesC) / 10.0)
}

// AmbientTemp returns the temperature reported by the CX34's ambient
// temperature sensor.
func (s *State) AmbientTemp() units.Temperature {
	deciDegreesC := s.registerValues[AmbientTemp]
	return units.FromCelsius(float64(deciDegreesC) / 10.0)
}

// InternalPumpSpeed returns the power setting of the variable-speed water pump inside
// of the CX34.
func (s *State) InternalPumpSpeed() units.PumpSpeed {
	return units.PumpSpeed(s.registerValues[InternalPumpSpeed])
}

// BoosterPumpSpeed returns the power setting of the variable-speed water pump
// external to the CX34. (not tested)
func (s *State) BoosterPumpSpeed() units.PumpSpeed {
	return units.PumpSpeed(s.registerValues[InternalPumpSpeed])
}

// ACVoltage returns the measured input AC Voltage value.
func (s *State) ACVoltage() units.Voltage {
	return units.Volt * units.Voltage(s.registerValues[InputACVoltage])
}

// ACCurrent returns the measured input AC Current value.
func (s *State) ACCurrent() units.Current {
	return units.Ampere * units.Current(s.registerValues[InputACCurrent]) / 10.0
}

// ApparentPower returns the measured input AC Current times the measured AC
// Voltage.
//
// I'm guessing there is no way to separate the real and reactive parts of the
// power value, so the actual power consumption is likely less than the returned
// value.
func (s *State) ApparentPower() units.Power {
	return units.PowerFromIV(s.ACCurrent(), s.ACVoltage())
}

// CompressorCurrent returns the "Compressor phase current value".
func (s *State) CompressorCurrent() units.Current {
	return units.Ampere * units.Current(s.registerValues[CompressorPhaseCurrent]) / 10.0
}

// InductorACCurrent returns the "inductor AC current value P15".
func (s *State) InductorACCurrent() units.Current {
	return units.Ampere * units.Current(s.registerValues[InductorACCurrent]) / 10.0
}

// UsefulHeatRate returns the amount of useful heat added or removed from the
// system per unit time. The value may be negative in the case of cooling.
func (s *State) UsefulHeatRate() units.Power {
	// H = delta T * specific heat of water * density of water
	massHeatedPerSec := waterDensity.TimesVolume(s.FlowRate().TimesDuration(time.Second))
	energyPerSec := constantSpecificHeat.TimesMassDeltaTemp(massHeatedPerSec, s.DeltaT())
	return units.Watt * units.Power(energyPerSec.Joules()) // 1 W = 1 Joule/sec
}

// UsefulHeatRateExplained returns the amount of useful heat added or removed from the
// system per unit time. The value may be negative in the case of cooling.
func (s *State) UsefulHeatRateExplained() string {
	// H = delta T * specific heat of water * density of water
	massHeatedPerSec := waterDensity.TimesVolume(s.FlowRate().TimesDuration(time.Second))
	energyPerSec := constantSpecificHeat.TimesMassDeltaTemp(massHeatedPerSec, s.DeltaT())
	return fmt.Sprintf("%.4fkg/s * %.1f°K * %.0fkJ/(kg * °K) = %.0fJ/s = %.2fkW",
		massHeatedPerSec.Kilograms(),
		s.DeltaT().Kelvin(),
		constantSpecificHeat.KilojoulesPerKilogramKelvin(),
		energyPerSec.Joules(),
		s.UsefulHeatRate().Kilowatts())
}

// MassFlowPerSecond returns the amount of water flowing through the heat pump per second.
func (s *State) MassFlowPerSecond() units.Mass {
	return waterDensity.TimesVolume(s.FlowRate().TimesDuration(time.Second))
}

// COP returns the coefficient of performance for the heat pump.
func (s *State) COP() (units.CoefficientOfPerformance, bool) {
	workRate := s.ApparentPower()
	if workRate == 0 {
		return 0, false
	}
	return units.CoefficientOfPerformance(s.UsefulHeatRate().Watts() / workRate.Watts()), true
}

// DeltaT returns the outlet temperature minust he inlet temperature
func (s *State) DeltaT() units.Temperature {
	return s.ACOutletWaterTemp() - s.ACInletWaterTemp()
}

/*

Table of registers with values that changed

| Register no. | Notes                       | Value |
|--------------|-----------------------------|-------|
| 143          | AC heating setpoint         | 39    |
| 144          | Domestic Hot Water setpoint | 51    |
| 200          | C0                          | 22    |
| 201          | C1                          | 700   |
| 202          | C2                          | 49    |
| 203          | C3                          | 8     |
| 204          | C4                          | 451   |
| 205          | C5                          | 453   |
| 206          | C6                          | 249   |
| 207          | C7                          | 245   |
| 208          | C8                          | 247   |
| 209          | C9                          | 60    |
| 213          | C13                         | 110   |
| 227          | C27                         | 43    |
| 229          | C29                         | 1     |
| 235          | C35                         | 0     |
| 237          | C37                         | 859   |
| 240          | C40                         | 98    |
| 241          | C41                         | 237   |
| 243          | C43                         | 0     |
| 245          | C45                         | 895   |
| 248          | C48                         | 6     |
| 250          | C50                         | 55    |
| 251          | C51                         | 15    |
| 255          | C55                         | 239   |
| 256          | C56                         | 55    |
| 257          | C57                         | 69    |
| 258          | C58                         | 389   |
| 260          | C60                         | 24    |
| 261          | C61                         | 14    |
| 280          | C80                         | 248   |
| 281          | C81                         | 399   |
| 282          | C82                         | 399   |
*/

// Known Register values.
const (
	TargetACHeatingModeTemp    Register = 143
	TargetDomesticHotWaterTemp Register = 144
	// See page 51 of https://www.chiltrix.com/documents/CX34-IOM-3.pdf
	OutPipeTemp               Register = 200
	CompressorDischargeTemp   Register = 201
	AmbientTemp               Register = 202
	SuctionTemp               Register = 203
	PlateHeatExchangerTemp    Register = 204
	ACOutletWaterTemp         Register = 205
	SolarTemp                 Register = 206
	CompressorCurrentValueP15 Register = 209 // 0.00-30.0A
	UsageSideWaterFlowRate    Register = 213 // tenths of a liter per minute
	P03Status                 Register = 214
	CompressorFrequency       Register = 227
	InternalPumpSpeed         Register = 248
	BoosterPumpSpeed          Register = 249
	InductorACCurrent         Register = 250
	InputACVoltage            Register = 255
	InputACCurrent            Register = 256
	CompressorPhaseCurrent    Register = 257

	WaterOutletSensorTemp1 Register = 204
	WaterOutletSensorTemp2 Register = 205
	WaterFlowRate          Register = 213
	WaterInletSensorTemp1  Register = 281
	WaterInletSensorTemp2  Register = 282
)

// Source: https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf
var registerNames = map[Register]string{
	143: "TargetACHeatingModeTemp",
	144: "TargetDomesticHotWaterTemp", // was: "Din7 AC Cooling Mode Switch",
	// Starting at 200, it's all the C parameters from the details screen.
	200: "ERR3.15",
	201: "EC Fan 1 Fault",
	202: "EC Fsn 2 Fault",
	203: "Heat Recovery Warning",
	204: "WaterOutletSensorTemp1",
	205: "WaterOutletSensorTemp2",
	213: "WaterFlowRate", // tenths of a liter / minute
	281: "WaterInletSensorTemp1",
	282: "WaterInletSensorTemp2",
}
