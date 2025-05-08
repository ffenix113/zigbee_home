package base

import (
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

// SensorType is type that will fetch only
type SensorType struct {
	Type string
}

// Sensor defines all information necessary about the attached sensor.
type Base struct {
	Type  string
	label string
	// Connection provides information about communication protocol.
	Connection map[string]string
}

func (b *Base) Label() string {
	return b.label
}

func (b *Base) SetLabel(label string) {
	if b.label != "" {
		panic("trying to set label second time")
	}

	b.label = label
}

// Includes is an optional include statements that will be written in main.cpp.
func (*Base) Includes() []string {
	return nil
}

func (*Base) Template() string {
	return ""
}

func (*Base) Clusters() cluster.Clusters {
	return nil
}

func (*Base) NeedsDevice() bool {
	return false
}

func (*Base) CPPComponentType() string {
	return ""
}

func (*Base) AppConfig() []appconfig.ConfigValue {
	return nil
}

func (*Base) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	return nil
}
