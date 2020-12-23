// Package tempsensor provides an API for reading DS18B20-based temperature
// sensors.
package tempsensor

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gonzojive/heatpump/mdtable"
)

const (
	guideURL = "https://www.circuitbasics.com/raspberry-pi-ds18b20-temperature-sensor-tutorial/#:~:text=The%20DS18B20%20temperature%20sensor%20is,accurate%20and%20take%20measurements%20quickly."
)

// Config is used to configure the paths used to search for temperature sensors.
type Config struct {
	// The root of the Raspberry Pi filesystem. When developing over sshfs, it
	// may be helpful to specify a different directory.
	Root string
}

func (c *Config) canonicalize() *Config {
	if c == nil {
		return DefaultConfig()
	}
	return c
}

func (c *Config) sensorListPath() string {
	return path.Join(c.Root, "sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
}

func (c *Config) w1SlavePath(id string) string {
	return path.Join(c.Root, "sys/bus/w1/devices", id, "w1_slave")
}

// DefaultConfig returns the config that's appropriate when the code is running
// on the Raspberry pi itself.
func DefaultConfig() *Config {
	return &Config{"/"}
}

// List returns a list of sensors
func List(c *Config) ([]*Sensor, error) {
	c = c.canonicalize()
	contents, err := ioutil.ReadFile(c.sensorListPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read list of temperature sensors from %s; follow the instructions at %s. to set up the temperature sensors on the Pi and try again: %w", c.sensorListPath(), guideURL, err)
	}
	ids := strings.Split(strings.TrimSpace(string(contents)), "\n")
	var sensors []*Sensor
	for _, id := range ids {
		if id == "" {
			// continue
		}
		sensors = append(sensors, &Sensor{id, c})
	}
	return sensors, nil
}

// Sensor is a temperature sensor.
type Sensor struct {
	id     string
	config *Config
}

// ID returns a unique identifier for the sensor. Each DS18B20 is provisioned
// with a unique hardware id that is used to identify it.
func (s *Sensor) ID() string {
	return s.id
}

// Read() reads the temperature of the sensor.
func (s *Sensor) Read() (Temperature, error) {
	filename := s.config.w1SlavePath(s.id)
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to read temperature file for device %q: %w", s.ID(), err)
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("problem reading temperature for device %q: contents of %s unexpected: %q", s.ID(), filename, string(contents))
	}
	if !strings.HasSuffix(lines[0], "YES") {
		return 0, fmt.Errorf("first line of %s is %q, want first line to end in YES", filename, lines[0])
	}

	matches := tempLineRegexp.FindStringSubmatch(lines[1])
	if len(matches) == 0 {
		return 0, fmt.Errorf("could not parse temperature of sensor %q: second line of %s is %q", s.ID(), filename, lines[1])
	}
	thousandths, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("could not parse temperature of sensor %q: %w", s.ID(), err)
	}

	return Temperature(thousandths), nil
}

var tempLineRegexp = regexp.MustCompile(`.* t=(.*)$`)

// Temperature is how this package deals with temperature values. The value is
// stored in thousandths of a degree centigrade.
type Temperature int32

// Celcius returns the temperature in celcius.
func (t Temperature) Celcius() float64 {
	return float64(t) / 1000
}

// Fahrenheit returns the temperature in degrees Fahrenheit.
func (t Temperature) Fahrenheit() float64 {
	return t.Celcius()*9/5 + 32
}

// String returns a human-readable representation of the temperature value.
func (t Temperature) String() string {
	return fmt.Sprintf("%.2f°C (%.2f°F)", t.Celcius(), t.Fahrenheit())
}

// DebugReport returns a string summary of all the temperature sensors. The
// second return value will be non-nil if an error was encountered.
func DebugReport(c *Config) (string, error) {
	sensors, err := List(c)
	if err != nil {
		err = fmt.Errorf("error reading sensors: %w", err)
		return fmt.Sprintf("error: %s", err.Error()), err
	}

	tb := &mdtable.Builder{}
	tb.SetHeader([]string{"Sensor ID", "Temperature"})

	for _, s := range sensors {
		t, err := s.Read()
		if err != nil {
			err = fmt.Errorf("error reading sensor %q: %w", s.id, err)
			return fmt.Sprintf("error: %s", err.Error()), err
		}
		tb.AddRow([]string{s.ID(), t.String()})
	}
	table := tb.Build()
	return fmt.Sprintf("%s\n", table), nil
}
