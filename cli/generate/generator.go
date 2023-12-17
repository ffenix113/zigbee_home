package generate

import (
	"fmt"
	"os"

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

func (g *Generator) Generate(workDir string, device *config.Device) error {
	// Write devicetree overlay (app.overlay)
	overlayFile, err := os.Create(workDir + "/app.overlay")
	if err != nil {
		return fmt.Errorf("create overlay file: %w", err)
	}

	defer overlayFile.Close()

	if err := g.DeviceTree.WriteTo(overlayFile); err != nil {
		return fmt.Errorf("write overlay: %w", err)
	}

	// Write app config (prj.conf)
	appConfigFile, err := os.Create(workDir + "/prj.conf")
	if err != nil {
		return fmt.Errorf("create app config file: %w", err)
	}

	defer appConfigFile.Close()

	if err := g.AppConfig.WriteTo(appConfigFile); err != nil {
		return fmt.Errorf("write app config: %w", err)
	}

	// Write app source
	srcDir := workDir + "/src"
	if err := os.Mkdir(srcDir, os.ModeDir|0o775); err != nil && !os.IsExist(err) {
		return fmt.Errorf("create src dir: %w", err)
	}

	if err := g.Source.WriteTo(srcDir, device); err != nil {
		return fmt.Errorf("write app src: %w", err)
	}

	return nil
}
