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
	OutPipeTemp                           Register = 200
	CompressorDischargeTemp               Register = 201
	AmbientTemp                           Register = 202
	SuctionTemp                           Register = 203
	PlateHeatExchangerTemp                Register = 204
	ACOutletWaterTemp                     Register = 205
	SolarTemp                             Register = 206
	CompressorCurrentValueP15             Register = 209 // 0.00-30.0A
	WaterFlowRate                         Register = 213 // tenths of a liter per minute
	P03Status                             Register = 214
	P04Status                             Register = 215
	P05Status                             Register = 216
	P06Status                             Register = 217
	P07Status                             Register = 218
	P08Status                             Register = 219 // 0= DHW valid, 1= DHW invalid 0=DHW valid, 1= DHW invalid
	P09Status                             Register = 220 // 0=Heating valid,	1= Heating invalid	AC heating valid= 0 valid, 	1= invalid
	P10Status                             Register = 221 // 0=cooling valid,	1=cooling invalid	0=cooling valid,	1=cooling invalid
	HighPressureSwitchStatus              Register = 222 // 1= on, 0= off 1= on, 0= off
	LowPressureSwitchStatus               Register = 223 // 1=on, 0= off 1=on, 0= off
	SecondHighPressureSwitchStatus        Register = 224 // 1=on, 0= off 1=on, 0= off
	InnerWaterFlowSwitch                  Register = 225 // 1=on, 0= off 1=on, 0= off
	CompressorFrequency                   Register = 227 // Displays the actual operating	frequency	Show actual frequency
	ThermalSwitchStatus                   Register = 228 // 1=on, 0= off 1=on, 0= off
	OutdoorFanMotor                       Register = 229 // 1= run, 0= stop 1=on, 0= off
	ElectricalValve1                      Register = 230 // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve2                      Register = 231 // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve3                      Register = 232 // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve4                      Register = 233 // 1= run, 0= stop 1= run, 0= stop
	C4WaterPump                           Register = 234 // 1= run, 0= stop 1= run, 0= stop
	C5WaterPump                           Register = 235 // 1= run, 0= stop 1= run, 0= stop
	C6waterPump                           Register = 236 // 1= run, 0= stop 1= run, 0= stop
	AccumulativeDaysAFterLastVirusKilling Register = 237 // The accumulative days after last	virus killing	0-99 (From the last complete	sterilization to the present,	cumulative number of days）	0-99 (from the last complete	sterilization to the present,	cumulative number of days)
	OutdoorModularTemp                    Register = 238 // -30~97℃ -30~97℃
	ExpansionValve1OpeningDegree          Register = 239 // 0~500 0~500
	ExpansionValve2OpeningDegree          Register = 240 // 0~500 0~500
	InnerPipeTemp                         Register = 241 // -30~97℃ -30~97℃
	HeatingMethod2TargetTemperature       Register = 242 // -30~97℃ -30~97℃
	IndoorTemperatureControlSwitch        Register = 243 // 1=on, 0= off 1=on, 0= off
	FanType                               Register = 244 // 0= AC fan, 1= EC fan 1,	2= EC fan 2	0= AC fan, 1= EC fan 1,	2= EC fan 2
	ECFanMotor1Speed                      Register = 245 // 0~3000 0~3000
	ECFanMotor2Speed                      Register = 246 //0~3000 0~3000
	WaterPumpTypes                        Register = 247 // 0= AC Water pump	1= EC Water pump	0= AC Water pump	1= EC Water pump
	InternalPumpSpeed                     Register = 248 // (C4) 1~10 （10 Show 100%） 1~10 (10 means 100%)
	BoosterPumpSpeed                      Register = 249 //1~10 （10 Show 100%） 1~10 (10 means 100%)
	InductorACCurrent                     Register = 250 //0~50A 0~50A
	DriverWorkingStatusValue              Register = 251 //Hexadecimal value Hexadecimal values
	CompressorShutDownCode                Register = 252 //Hexadecimal value Hexadecimal values
	DriverAllowedHighestFrequency         Register = 253 //30-120Hz 30-120Hz
	ReduceFrequencyTemperature            Register = 254 //setting	55~200℃ 55~200℃
	InputACVoltage                        Register = 255 //0~550V 0~550V
	InputACCurrent                        Register = 256 //0~50A（IPM test） 0~50A（IPM Check）
	CompressorPhaseCurrent                Register = 257 //0~50A（IPM test） 0~50A（IPM Check）
	BusLineVoltage                        Register = 258 //0~750V 0~750V
	FanShutdownCode                       Register = 259 // Hexadecimal value Hexadecimal values
	IPMTemp                               Register = 260 //55~200℃ 55~200℃
	CompressorTotalRunningTime            Register = 261 //	Will reset after power cycle	0~65000 0~65000 hour

	// Inferred values.
	WaterInletSensorTemp1 Register = 281
	WaterInletSensorTemp2 Register = 282
	CurrentFaultCode      Register = 284 // Set to 32 when I get a P5 error, not sure about other faults.
)

// Source: https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf
var registerNames = map[Register]string{
	143: "TargetACHeatingModeTemp",
	144: "TargetDomesticHotWaterTemp", // was: "Din7 AC Cooling Mode Switch",
	// Starting at 200, it's all the C parameters from the details screen.
	WaterInletSensorTemp1: "WaterInletSensorTemp1",
	WaterInletSensorTemp2: "WaterInletSensorTemp2",

	OutPipeTemp:                           "OutPipeTemp",
	CompressorDischargeTemp:               "CompressorDischargeTemp",
	AmbientTemp:                           "AmbientTemp",
	SuctionTemp:                           "SuctionTemp",
	PlateHeatExchangerTemp:                "PlateHeatExchangerTemp",
	ACOutletWaterTemp:                     "ACOutletWaterTemp",
	SolarTemp:                             "SolarTemp",
	CompressorCurrentValueP15:             "CompressorCurrentValueP15", // 0.00-30.0A
	WaterFlowRate:                         "WaterFlowRate",             // tenths of a liter per minute
	P03Status:                             "P03Status",
	P04Status:                             "P04Status",
	P05Status:                             "P05Status",
	P06Status:                             "P06Status",
	P07Status:                             "P07Status",
	P08Status:                             "P08Status",                             // 0= DHW valid, 1= DHW invalid 0=DHW valid, 1= DHW invalid
	P09Status:                             "P09Status",                             // 0=Heating valid,	1= Heating invalid	AC heating valid= 0 valid, 	1= invalid
	P10Status:                             "P10Status",                             // 0=cooling valid,	1=cooling invalid	0=cooling valid,	1=cooling invalid
	HighPressureSwitchStatus:              "HighPressureSwitchStatus",              // 1= on, 0= off 1= on, 0= off
	LowPressureSwitchStatus:               "LowPressureSwitchStatus",               // 1=on, 0= off 1=on, 0= off
	SecondHighPressureSwitchStatus:        "SecondHighPressureSwitchStatus",        // 1=on, 0= off 1=on, 0= off
	InnerWaterFlowSwitch:                  "InnerWaterFlowSwitch",                  // 1=on, 0= off 1=on, 0= off
	CompressorFrequency:                   "CompressorFrequency",                   // Displays the actual operating	frequency	Show actual frequency
	ThermalSwitchStatus:                   "ThermalSwitchStatus",                   // 1=on, 0= off 1=on, 0= off
	OutdoorFanMotor:                       "OutdoorFanMotor",                       // 1= run, 0= stop 1=on, 0= off
	ElectricalValve1:                      "ElectricalValve1",                      // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve2:                      "ElectricalValve2",                      // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve3:                      "ElectricalValve3",                      // 1= run, 0= stop 1= run, 0= stop
	ElectricalValve4:                      "ElectricalValve4",                      // 1= run, 0= stop 1= run, 0= stop
	C4WaterPump:                           "C4WaterPump",                           // 1= run, 0= stop 1= run, 0= stop
	C5WaterPump:                           "C5WaterPump",                           // 1= run, 0= stop 1= run, 0= stop
	C6waterPump:                           "C6waterPump",                           // 1= run, 0= stop 1= run, 0= stop
	AccumulativeDaysAFterLastVirusKilling: "AccumulativeDaysAFterLastVirusKilling", // The accumulative days after last	virus killing	0-99 (From the last complete	sterilization to the present,	cumulative number of days）	0-99 (from the last complete	sterilization to the present,	cumulative number of days)
	OutdoorModularTemp:                    "OutdoorModularTemp",                    // -30~97℃ -30~97℃
	ExpansionValve1OpeningDegree:          "ExpansionValve1OpeningDegree",          // 0~500 0~500
	ExpansionValve2OpeningDegree:          "ExpansionValve2OpeningDegree",          // 0~500 0~500
	InnerPipeTemp:                         "InnerPipeTemp",                         // -30~97℃ -30~97℃
	HeatingMethod2TargetTemperature:       "HeatingMethod2TargetTemperature",       // -30~97℃ -30~97℃
	IndoorTemperatureControlSwitch:        "IndoorTemperatureControlSwitch",        // 1=on, 0= off 1=on, 0= off
	FanType:                               "FanType",                               // 0= AC fan, 1= EC fan 1,	2= EC fan 2	0= AC fan, 1= EC fan 1,	2= EC fan 2
	ECFanMotor1Speed:                      "ECFanMotor1Speed",                      // 0~3000 0~3000
	ECFanMotor2Speed:                      "ECFanMotor2Speed",                      //0~3000 0~3000
	WaterPumpTypes:                        "WaterPumpTypes",                        // 0= AC Water pump	1= EC Water pump	0= AC Water pump	1= EC Water pump
	InternalPumpSpeed:                     "InternalPumpSpeed",                     // (C4) 1~10 （10 Show 100%） 1~10 (10 means 100%)
	BoosterPumpSpeed:                      "BoosterPumpSpeed",                      //1~10 （10 Show 100%） 1~10 (10 means 100%)
	InductorACCurrent:                     "InductorACCurrent",                     //0~50A 0~50A
	DriverWorkingStatusValue:              "DriverWorkingStatusValue",              //Hexadecimal value Hexadecimal values
	CompressorShutDownCode:                "CompressorShutDownCode",                //Hexadecimal value Hexadecimal values
	DriverAllowedHighestFrequency:         "DriverAllowedHighestFrequency",         //30-120Hz 30-120Hz
	ReduceFrequencyTemperature:            "ReduceFrequencyTemperature",            //setting	55~200℃ 55~200℃
	InputACVoltage:                        "InputACVoltage",                        //0~550V 0~550V
	InputACCurrent:                        "InputACCurrent",                        //0~50A（IPM test） 0~50A（IPM Check）
	CompressorPhaseCurrent:                "CompressorPhaseCurrent",                //0~50A（IPM test） 0~50A（IPM Check）
	BusLineVoltage:                        "BusLineVoltage",                        //0~750V 0~750V
	FanShutdownCode:                       "FanShutdownCode",                       // Hexadecimal value Hexadecimal values
	IPMTemp:                               "IPMTemp",                               //55~200℃ 55~200℃
	CompressorTotalRunningTime:            "CompressorTotalRunningTime",            //	Will reset after power cycle	0~65000 0~65000 hour

	CurrentFaultCode: "Fault Code?", // Set to 32 when I get a P5 error code.
}
