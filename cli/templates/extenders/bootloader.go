package extenders

import (
	"github.com/ffenix113/zigbee_home/cli/types/board"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
)

func NewBootloaderConfig(config *board.Bootloader) generator.SimpleExtender {
	return generator.SimpleExtender{
		Name: "Bootloader configuration",
		FilesToWrite: []generator.WriteFile{
			{
				FileName:          "../pm_static.yml",
				TemplateName:      "pm_static.yml",
				AdditionalContext: config.PM,
			},
		},
		Config: config.Config,
	}
}
