package extenders

import (
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = Sensor{}
var _ appconfig.Provider = Sensor{}

type Sensor struct{}

func NewSensor() Sensor {
	return Sensor{}
}

// AppConfig implements appconfig.Provider.
func (Sensor) AppConfig() []appconfig.ConfigValue {
	return []appconfig.ConfigValue{
		appconfig.CONFIG_SENSOR.Required(appconfig.Yes),
	}
}

// Includes implements templates.Extender.
func (Sensor) Includes() []string {
	return []string{"zephyr/drivers/sensor.h", "zbhome_sensor.hpp"}
}

// Template implements templates.Extender.
func (Sensor) Template() string {
	return ""
}

// WriteFiles implements templates.Extender.
func (Sensor) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "zbhome_sensor.hpp",
			TemplateName: "zbhome_sensor.hpp",
		},
		{
			FileName:     "zbhome_sensor.cpp",
			TemplateName: "zbhome_sensor.cpp",
		},
	}
}

func (Sensor) ZephyrModules() []string {
	return nil
}
