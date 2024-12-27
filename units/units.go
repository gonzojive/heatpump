// Package units contains float64 wrappers with unit semantics relevant to the
// CX34 heat pump (flow rate, temperature, etc).
package units

import (
	"fmt"
	"time"

	underlying "github.com/martinlindhe/unit"
)

// Temperature contains a temperature value in kelvin.
type Temperature = underlying.Temperature

// FlowRate is measured in units of volume over duration, such as liters per
// minute. The materialized value is liters/minute.
type FlowRate float64

// Flow rate values.
const (
	LiterPerMinute FlowRate = FlowRate(Liter)
)

// LitersPerMinute returns the flow rate in liters per minute as a float64.
func (v FlowRate) LitersPerMinute() float64 {
	return Volume(v).Liters()
}

// LitersPerSecond returns the flow rate in liters per minute as a float64.
func (v FlowRate) LitersPerSecond() float64 {
	return v.LitersPerMinute() / 60
}

// TimesDuration returns the volume of flow over the course of a given duration.
func (v FlowRate) TimesDuration(dur time.Duration) Volume {
	return Liter * Volume(dur.Minutes()*v.LitersPerMinute())
}

// Scale scales the value by a floating point multiplier.
func (v FlowRate) Scale(f float64) FlowRate {
	return v * FlowRate(f)
}

// FromCelsius returns a temperature from a value in degrees celcius.
func FromCelsius(t float64) Temperature {
	return underlying.FromCelsius(t)
}

// Current represents electric current.
type Current = underlying.ElectricCurrent

const Ampere = underlying.Ampere

// Voltage represents electric Voltage.
type Voltage = underlying.Voltage

const Volt = underlying.Volt

// Power represents energy over time.
type Power = underlying.Power

// Power constants.
const (
	Watt     Power = underlying.Watt
	Kilowatt Power = underlying.Kilowatt
)

// PowerFromIV multiples current and voltage to get power.
func PowerFromIV(i Current, v Voltage) Power {
	return underlying.Watt * underlying.Power(i.Amperes()*v.Volts())
}

// PumpSpeed is the speed of the pump. The value should be between 0 (for off)
// and 10 (for 100% on).
type PumpSpeed uint8

// String returns the pump speed as a fraction of full power.
func (s PumpSpeed) String() string {
	if s == 0 {
		return "0 (off)"
	}
	return fmt.Sprintf("%d/10", int(s))
}

// AsFractionOfMaxSpeed returns the pump speed divided by the max pump speed.
func (s PumpSpeed) AsFractionOfMaxSpeed() float32 {
	return float32(s) / 10
}

// CoefficientOfPerformance is the ratio of useful heat supplied or removed from
// the water stream divided by the amount of electrical energy energy being used
// by the heat pump.
type CoefficientOfPerformance float64

// Float64 returns the COP as a floating point number.
func (s CoefficientOfPerformance) Float64() float64 {
	return float64(s)
}

// SpecificHeat is measured in kJ/(kg * K).
type SpecificHeat float64

// KilojoulePerKilogramKelvin is the standard unit used for recording specific
// heat in tables.
const KilojoulePerKilogramKelvin SpecificHeat = 1

// Scale scales the specific heat value by a multiplier.
func (sh SpecificHeat) Scale(f float64) SpecificHeat {
	return SpecificHeat(f) * sh
}

// KilojoulesPerKilogramKelvin returns the specific heat in kJ/(kg * Kelvin).
func (sh SpecificHeat) KilojoulesPerKilogramKelvin() float64 {
	return float64(sh)
}

// TimesMassDeltaTemp multiplies the specific heat value by mass and delta T to
// arrive at an energy value.
func (sh SpecificHeat) TimesMassDeltaTemp(m Mass, t Temperature) Energy {
	return underlying.Kilojoule * underlying.Energy(float64(sh)*m.Kilograms()*t.Kelvin())
}

// Mass is a floating point mass value.
type Mass = underlying.Mass

// Mass values.
const (
	Kilogram = underlying.Kilogram
)

// Energy is a floating point energy value.
type Energy = underlying.Energy

// Volume is a floating point volume value.
type Volume = underlying.Volume

// Volume values.
const (
	CubicMeter = underlying.CubicMeter
	Liter      = underlying.Liter
)

// Density is a the density of a material (mass/volume).
//
// The materialized value is kg/m^3.
type Density float64

// DensityFromRatio returns a density value that is the ratio of a given mass and volume.
func DensityFromRatio(m Mass, v Volume) Density {
	return Density(m.Kilograms() / v.CubicMeters())
}

// TimesVolume returns the mass taken up by a given volume of material with
// density p.
func (p Density) TimesVolume(v Volume) Mass {
	return Kilogram * Mass(Mass(p).Kilograms()*v.CubicMeters())
}
