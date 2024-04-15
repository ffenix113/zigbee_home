package sensor

import (
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

var _ Sensor = (*Simple)(nil)

type Simple struct {
	SensorName       string
	SensorLabel      string
	SensorTemplate   string
	SensorClusters   cluster.Clusters
	SensorAppConfig  []appconfig.ConfigValue
	SensorAppOverlay func(*devicetree.DeviceTree) error
}

func (s *Simple) String() string {
	return s.SensorName
}

func (s *Simple) Label() string {
	return s.SensorLabel
}

func (s *Simple) Template() string {
	return s.SensorTemplate
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
