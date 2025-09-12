package base

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
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

func (*SoilMoistureADC) CPPComponentType() string {
	return "SoilMoistureADC"
}

func (o *SoilMoistureADC) Clusters() cluster.Clusters {
	return []cluster.Cluster{
		// Hardcoded, as we don't configure this values.
		cluster.NewSoilMoisture(0, 100),
	}
}

func (*SoilMoistureADC) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "types_soil_moisture_adc.hpp",
			TemplateName: "types_soil_moisture_adc.hpp",
		},
		{
			FileName:     "types_soil_moisture_adc.cpp",
			TemplateName: "types_soil_moisture_adc.cpp",
		},
		{
			FileName:     "zbhome_sensor.cpp",
			TemplateName: "zbhome_sensor.cpp",
		},
		{
			FileName:     "zbhome_sensor.hpp",
			TemplateName: "zbhome_sensor.hpp",
		},
	}
}

func (*SoilMoistureADC) Includes() []string {
	return []string{
		"zephyr/drivers/adc.h",
		"types_soil_moisture_adc.hpp",
	}
}

func (*SoilMoistureADC) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_ADC").Required(appconfig.Yes),
	}
}

func (o *SoilMoistureADC) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	if err := o.ADCPin.AttachSelf(overlay); err != nil {
		return fmt.Errorf("attach adc pin: %w", err)
	}

	dtPin := devicetree.NewButton(o.ADCPin.Pin)
	return dtPin.AttachSelf(overlay)
}
