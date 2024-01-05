package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/sensor"
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/sensor/bosch"
)

var knownSensors = map[string]Sensor{
	// Generic
	"on_off": &base.OnOff{},

	// Specific devices
	"bme280":             &bosch.BME280{},
	"device_temperature": &sensor.DeviceTemperature{},
}
