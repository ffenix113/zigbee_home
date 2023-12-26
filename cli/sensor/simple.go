package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/generate/appconfig"
	"github.com/ffenix113/zigbee_home/cli/generate/devicetree"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type Simple struct {
	SensorName       string
	SensorClusters   cluster.Clusters
	SensorAppConfig  []appconfig.ConfigValue
	SensorAppOverlay func(*devicetree.DeviceTree) error
}

func (s *Simple) Name() string {
	return s.SensorName
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
