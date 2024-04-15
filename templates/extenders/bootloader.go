package extenders

import (
	"github.com/ffenix113/zigbee_home/types/board"
	"github.com/ffenix113/zigbee_home/types/generator"
)

func NewBootloaderConfig(config *board.Bootloader) generator.SimpleExtender {
	extender := generator.SimpleExtender{
		Name:   "Bootloader configuration",
		Config: config.Config,
	}

	if config.PM != nil {
		extender.FilesToWrite = []generator.WriteFile{
			{
				FileName:          "../pm_static.yml",
				TemplateName:      "pm_static.yml",
				AdditionalContext: config.PM,
			},
		}
	}

	return extender
}
