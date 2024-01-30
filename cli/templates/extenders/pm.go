package extenders

import (
	"github.com/ffenix113/zigbee_home/cli/types/board"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

func NewStaicPM(config board.PMConfig) generator.SimpleExtender {
	return generator.SimpleExtender{
		Name: "Static PM",
		FilesToWrite: []generator.WriteFile{
			{
				FileName:          "../pm_static.yml",
				TemplateName:      "pm_static.yml",
				AdditionalContext: config,
			},
		},
	}
}
