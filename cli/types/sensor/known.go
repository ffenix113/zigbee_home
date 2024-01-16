package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/sensor"
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/sensor/bosch"
	"github.com/ffenix113/zigbee_home/cli/sensor/sensirion"
)

var knownSensors = map[string]Sensor{
	// Generic
	"on_off": &base.OnOff{},

	// Specific devices

	"device_temperature": &sensor.DeviceTemperature{},

	// Bosch
	"bme280": &bosch.BME280{},
	// This is a clone of bme280, with different overlay name
	// FIXME: It does not yet support IAQ measurements,
	// and does not expose resistance to Zigbee.
	"bme680": bosch.NewBME680(),

	// Sensirion
	"scd4x": &sensirion.SCD4X{},
}
