package base

import (
	"github.com/ffenix113/zigbee_home/cli/generate/appconfig"
	"github.com/ffenix113/zigbee_home/cli/generate/devicetree"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

// SensorType is type that will fetch only
type SensorType struct {
	Type string
}

// Sensor defines all information necessary about the attached sensor.
type Base struct {
	Type string `yaml:"type"`
	// Connection provides information about communication protocol.
	Connection map[string]string
}

func (*Base) Clusters() cluster.Clusters {
	return nil
}

func (*Base) AppConfig() []appconfig.ConfigValue {
	return nil
}

func (*Base) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	return nil
}
