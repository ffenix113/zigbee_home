package extenders

import "github.com/ffenix113/zigbee_home/types/generator"

type SCD4X struct {
	generator.SimpleExtender
}

func NewSCD4X() SCD4X {
	return SCD4X{
		SimpleExtender: generator.SimpleExtender{
			ZephyrModuleNames: []string{"scd4x"},
		},
	}
}
