package bosch

import (
	"github.com/ffenix113/zigbee_home/cli/generate/appconfig"
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type BME280 struct {
	*base.Base `yaml:",inline"`
}

func (BME280) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.Temperature{},
		// TODO: humidity, pressure
	}
}

func (BME280) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.CONFIG_BME280.Required(appconfig.Yes),
	}
}
