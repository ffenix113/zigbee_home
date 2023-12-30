package sensor

import (
	"github.com/ffenix113/zigbee_home/cli/sensor/base"
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

var _ appconfig.Provider = &InternalTemperature{}

type InternalTemperature struct {
	*base.Base `yaml:",inline"`
}

func (*InternalTemperature) String() string {
	return "internal temperature"
}

func (*InternalTemperature) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_NRFX_TEMP").Required(appconfig.Yes),
	}
}

func (*InternalTemperature) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.DeviceTemperature{},
	}
}

func (*InternalTemperature) Template() string {
	return "sensors/internal_temperature"
}

func (*InternalTemperature) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewNrfxTemp(),
	}
}
