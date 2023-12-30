package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/sensor"
	"github.com/ffenix113/zigbee_home/cli/sensor/bosch"
)

var knownSensors = map[string]Sensor{
	"bme280":               &bosch.BME280{},
	"internal_temperature": &sensor.InternalTemperature{},
}
