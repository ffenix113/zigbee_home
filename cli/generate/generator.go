package generate

import (
	"fmt"
	"os"

	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
	"github.com/ffenix113/zigbee_home/cli/types/devicetree"
	"github.com/ffenix113/zigbee_home/cli/types/generator"
	"github.com/ffenix113/zigbee_home/cli/types/sensor"
	"github.com/ffenix113/zigbee_home/cli/types/source"
)

// TODO: Define some "modifier"/"injector" interface
// 	that will allow to update generator state.
// 	It is important that this interface be able to add
// 	new files to generated source output!
//
// 	This is useful for common functionality for
// 	sensors and peripherals.
// 	I.e. to enable some sensors it is neeeded:
// 	 - update proj.conf to add CONFIG_SENSOR=y
// 	 - - maybe also add CONFIG_I2C=y
// 	 - add include <zephyr/drivers/sensor.h>
//
// 	Sensors should also be able to tell that
// 	they want to use such "modifier"/"injector"

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
	var providedExtenders []generator.Extender
	uniqueExtenders := map[string]struct{}{}

	for _, deviceSensor := range device.Sensors {
		if withExtender, ok := deviceSensor.(sensor.WithExtenders); ok {
			sensorExtenders := withExtender.Extenders()
			cleanedExtenders := sensorExtenders[:0]

			for _, extender := range sensorExtenders {
				extenderName := generator.ExtenderName(extender)

				if _, ok := uniqueExtenders[extenderName]; !ok {
					cleanedExtenders = append(cleanedExtenders, extender)
					uniqueExtenders[extenderName] = struct{}{}
				}
			}

			providedExtenders = append(providedExtenders, cleanedExtenders...)
		}
	}

	if device.Board.I2C != nil {
		providedExtenders = append(providedExtenders, extenders.NewI2C(device.Board.I2C...))
	}

	if device.Board.DebugLog {
		providedExtenders = append(providedExtenders,
			extenders.NewUSBUARTLog(),
		)
	}

	// Write devicetree overlay (app.overlay)
	if err := updateDeviceTree(device, g.DeviceTree, providedExtenders); err != nil {
		return fmt.Errorf("update overlay: %w", err)
	}

	overlayFile, err := os.Create(workDir + "/app.overlay")
	if err != nil {
		return fmt.Errorf("create overlay file: %w", err)
	}

	defer overlayFile.Close()

	if err := g.DeviceTree.WriteTo(overlayFile); err != nil {
		return fmt.Errorf("write overlay: %w", err)
	}

	// Write app config (prj.conf)
	if err := updateAppConfig(device, g.AppConfig, providedExtenders); err != nil {
		return fmt.Errorf("update app config: %w", err)
	}

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

	if err := g.Source.WriteTo(srcDir, device, providedExtenders); err != nil {
		return fmt.Errorf("write app src: %w", err)
	}

	return nil
}

func updateDeviceTree(device *config.Device, deviceTree *devicetree.DeviceTree, extenders []generator.Extender) error {
	for _, extender := range extenders {
		if withOverlayApplier, ok := extender.(devicetree.Applier); ok {
			if err := withOverlayApplier.ApplyOverlay(deviceTree); err != nil {
				return fmt.Errorf("apply extender: %w", err)
			}
		}
	}

	for _, sensor := range device.Sensors {
		if err := sensor.ApplyOverlay(deviceTree); err != nil {
			return fmt.Errorf("applying sensor %q to device tree: %w", sensor, err)
		}
	}

	return nil
}

func updateAppConfig(device *config.Device, appConfig *appconfig.AppConfig, extenders []generator.Extender) error {
	for _, extender := range extenders {
		if withAppConfigProvider, ok := extender.(appconfig.Provider); ok {
			appConfig.AddValue(withAppConfigProvider.AppConfig()...)
		}
	}

	for _, sensor := range device.Sensors {
		appConfig.AddValue(sensor.AppConfig()...)
	}

	return nil
}
