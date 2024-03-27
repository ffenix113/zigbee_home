package base

import (
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type PowerConfiguration struct {
	*Base                      `yaml:",inline"`
	cluster.PowerConfiguration `yaml:",inline"`
	ADCPin                     devicetree.ADCPin `yaml:"adc_pin"`
}

func (*PowerConfiguration) String() string {
	return "PowerConfiguration"
}

func (*PowerConfiguration) Template() string {
	return "sensors/common"
}

func (o *PowerConfiguration) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		o.PowerConfiguration,
	}
}

func (*PowerConfiguration) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_ADC").Required(appconfig.Yes),
	}
}

func (o *PowerConfiguration) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	dtPin := devicetree.NewButton(o.ADCPin.Pin)
	return dtPin.AttachSelf(overlay)
}

func (c *PowerConfiguration) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewADC(c.ADCPin),
	}
}
