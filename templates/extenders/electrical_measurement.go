package extenders

import (
	"github.com/ffenix113/zigbee_home/types/generator"
)

var _ generator.Extender = ElectricalMeasurement{}

type ElectricalMeasurement struct {
	generator.SimpleExtender
}

func (l ElectricalMeasurement) WriteFiles() []generator.WriteFile {
	return []generator.WriteFile{
		{
			FileName:     "cluster_electrical_measurement.hpp",
			TemplateName: "cluster_electrical_measurement.hpp",
		},
		{
			FileName:     "cluster_electrical_measurement.cpp",
			TemplateName: "cluster_electrical_measurement.cpp",
		},
	}
}

func (l ElectricalMeasurement) Includes() []string {
	return []string{
		"cluster_electrical_measurement.hpp",
		"zbhome_sensor.hpp",
	}
}
