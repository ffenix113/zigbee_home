package generate

import (
	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/generate/appconfig"
	"github.com/ffenix113/zigbee_home/cli/generate/devicetree"
	"github.com/ffenix113/zigbee_home/cli/generate/source"
)

type Generator struct {
	AppConfig  *appconfig.AppConfig
	DeviceTree *devicetree.DeviceTree
	Source     *source.Source
}

func NewGenerator(device config.Device) *Generator {
	return &Generator{
		AppConfig:  appconfig.NewDefaultAppConfig(),
		DeviceTree: devicetree.NewDeviceTree(),
		Source:     source.NewSource(),
	}
}

func (g *Generator) Generate(workDir string) error {
	return nil
}
