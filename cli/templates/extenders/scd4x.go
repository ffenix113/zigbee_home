package extenders

import "github.com/ffenix113/zigbee_home/cli/types/generator"

type SCD4X struct {
	generator.SimpleExtender
}

func NewSCD4X() SCD4X {
	return SCD4X{
		SimpleExtender: generator.SimpleExtender{
			ZephyrModuleNames: []string{"scd4x"},
			FilesToWrite: []generator.WriteFile{
				{
					FileName:     "scd4x.h",
					TemplateName: "sensors/sensirion/scd4x.h",
				},
				{
					FileName:     "scd4x.c",
					TemplateName: "sensors/sensirion/scd4x.c",
				},
			},
		},
	}
}
