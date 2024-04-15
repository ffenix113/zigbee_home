package extenders

import "github.com/ffenix113/zigbee_home/types/generator"

func NewNrfxTemp() generator.SimpleExtender {
	return generator.SimpleExtender{
		Name:           "NRFX Temp",
		IncludeHeaders: []string{"nrfx_temp.h"},
	}
}
