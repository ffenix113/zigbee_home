package base

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/zcl/cluster"
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
	return "sensors/power_config"
}

func (*PowerConfiguration) CPPComponentType() string {
	return "PowerConfig"
}

func (o *PowerConfiguration) Clusters() cluster.Clusters {
	clusterConfig := o.PowerConfiguration
	clusterConfig.BatteryRatedVoltage /= 100
	clusterConfig.BatteryVoltageMinThreshold /= 100
	return []cluster.Cluster{
		clusterConfig,
	}
}

func (*PowerConfiguration) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "types_power_config.hpp",
			TemplateName: "types_power_config.hpp",
		},
		{
			FileName:     "types_power_config.cpp",
			TemplateName: "types_power_config.cpp",
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

func (*PowerConfiguration) Includes() []string {
	return []string{
		"zephyr/drivers/adc.h",
		"types_power_config.hpp",
	}
}

func (*PowerConfiguration) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.NewValue("CONFIG_ADC").Required(appconfig.Yes),
	}
}

func (o *PowerConfiguration) ApplyOverlay(overlay *devicetree.DeviceTree) error {
	if err := o.ADCPin.AttachSelf(overlay); err != nil {
		return fmt.Errorf("attach adc pin: %w", err)
	}

	return nil
}
