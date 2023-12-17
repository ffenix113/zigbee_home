package sensor

import "github.com/ffenix113/zigbee_home/cli/sensor/bosch"

var knownSensors = map[string]any{
	"bme280": &bosch.BME280{},
}
