package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

var _ Sensor = (*Simple)(nil)

type Simple struct {
	SensorName           string
	SensorTemplatePrefix string
	SensorClusters       cluster.Clusters
	SensorAppConfig      []appconfig.ConfigValue
	SensorAppOverlay     func(*devicetree.DeviceTree) error
}

func (s *Simple) String() string {
	return s.SensorName
}

func (s *Simple) TemplatePrefix() string {
	return s.SensorTemplatePrefix
}

func (s *Simple) Clusters() cluster.Clusters {
	return s.SensorClusters
}

func (s *Simple) AppConfig() []appconfig.ConfigValue {
	return s.SensorAppConfig
}

func (s *Simple) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	return s.SensorAppOverlay(overlay)
}
