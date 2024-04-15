package sensor

import (
	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
)

var _ appconfig.Provider = &DeviceTemperature{}

type DeviceTemperature struct {
	*base.Base `yaml:",inline"`
}

func (*DeviceTemperature) String() string {
	return "device temperature"
}

func (*DeviceTemperature) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_NRFX_TEMP").Required(appconfig.Yes),
	}
}

func (*DeviceTemperature) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		cluster.DeviceTemperature{},
	}
}

func (*DeviceTemperature) Template() string {
	return "sensors/device_temperature"
}

func (*DeviceTemperature) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewNrfxTemp(),
	}
}
