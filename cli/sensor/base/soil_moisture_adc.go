package base

import (
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/zcl/cluster"
)

type SoilMoistureADC struct {
	*Base         `yaml:",inline"`
	MinMoistureMv uint16            `yaml:"min_moisture_mv"`
	MaxMoistureMv uint16            `yaml:"max_moisture_mv"`
	ADCPin        devicetree.ADCPin `yaml:"adc_pin"`
}

func (*SoilMoistureADC) String() string {
	return "SoilMoistureADC"
}

func (*SoilMoistureADC) Template() string {
	return "sensors/soil_moisture_adc"
}

func (o *SoilMoistureADC) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		// Hardcoded, as we don't configure this values.
		cluster.NewSoilMoisture(0, 100),
	}
}

func (*SoilMoistureADC) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_ADC").Required(appconfig.Yes),
	}
}

func (o *SoilMoistureADC) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	dtPin := devicetree.NewButton(o.ADCPin.Pin)
	return dtPin.AttachSelf(overlay)
}

func (c *SoilMoistureADC) Extenders() []generator.Extender {
	return []generator.Extender{
		extenders.NewADC(c.ADCPin),
	}
}
