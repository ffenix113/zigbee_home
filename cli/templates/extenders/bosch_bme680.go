package extenders

import "github.com/ffenix113/zigbee_home/cli/types/generator"

func NewBoschBME680() generator.SimpleExtender {
	return generator.SimpleExtender{
		Name:           "Bosch_BME680",
		IncludeHeaders: []string{"bosch_bme680.h"},
		FilesToWrite: []generator.WriteFile{
			{
				FileName:     "bosch_bme680.h",
				TemplateName: "sensors/bosch/bme680.h",
			},
			{
				FileName:     "bosch_bme680.c",
				TemplateName: "sensors/bosch/bme680.c",
			},
		},
	}
}
