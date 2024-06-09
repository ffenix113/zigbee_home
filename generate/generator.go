package generate

import (
	"fmt"
	"log"
	"os"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types/appconfig"
	"github.com/ffenix113/zigbee_home/types/board"
	"github.com/ffenix113/zigbee_home/types/devicetree"
	"github.com/ffenix113/zigbee_home/types/generator"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/ffenix113/zigbee_home/types/source"
)

type Generator struct {
	AppConfig  *appconfig.AppConfig
	DeviceTree *devicetree.DeviceTree
	Source     *source.Source
}

func NewGenerator(device *config.Device) (*Generator, error) {
	appConfig, err := appconfig.NewDefaultAppConfig(
		appconfig.DefaultAppConfigOptions{
			IsRouter:       device.Board.IsRouter,
			ZigbeeChannels: device.General.ZigbeeChannels,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("default app config: %w", err)
	}

	return &Generator{
		AppConfig:  appConfig,
		DeviceTree: devicetree.NewDeviceTree(),
		Source:     source.NewSource(),
	}, nil
}

func (g *Generator) Generate(workDir string, device *config.Device) error {
	providedExtenders, err := getExtenders(device)
	if err != nil {
		return fmt.Errorf("get extenders: %w", err)
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

func getExtenders(device *config.Device) ([]generator.Extender, error) {
	var providedExtenders []generator.Extender

	uniqueExtenders := map[string]struct{}{}

	forcedBootloader := device.Board.Bootloader != nil

	var bootloaderName string
	if forcedBootloader {
		bootloaderName = *device.Board.Bootloader
	}

	bootloaderConfig, bootloaderName := getBootloaderConfig(device.General.Board, bootloaderName)
	log.Printf("Device: %q, selected bootloader: %q, forced bootloader: %t\n",
		device.General.Board,
		bootloaderName,
		forcedBootloader)

	if forcedBootloader && bootloaderConfig == nil {
		return nil, fmt.Errorf("Bootloader %q was forced, but is not found in known bootloaders", *device.Board.Bootloader)
	}

	if bootloaderConfig != nil {
		providedExtenders = append(providedExtenders, extenders.NewBootloaderConfig(bootloaderConfig))
	}

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

	if device.Board.Debug != nil && device.Board.Debug.Enabled {
		providedExtenders = append(providedExtenders, extenders.NewDebugUARTLog(*device.Board.Debug))

		if device.Board.Debug.Console == extenders.DebugConsoleUSB {
			providedExtenders = append(providedExtenders, extenders.NewUSBUART())
		}
	}

	if len(device.Board.UART) != 0 {
		providedExtenders = append(providedExtenders, extenders.NewUART(device.Board.UART...))
	}

	if len(device.Board.LEDs) != 0 {
		providedExtenders = append(providedExtenders, extenders.NewLEDs(device.Board.LEDs.ToPins()...))
	}

	if len(device.Board.Buttons) != 0 {
		providedExtenders = append(providedExtenders, extenders.NewButtons(device.Board.Buttons.ToPins()...))
	}

	return providedExtenders, nil
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

func getBootloaderConfig(boardName, bootloader string) (*board.Bootloader, string) {
	if bootloader != "" {
		// Force using specific bootloader, even if it is unknown.
		// Unknown bootloaders will result in no additional configuration.
		return board.BootloaderConfig(bootloader), bootloader
	}

	return board.BootloaderConfigFromBoard(boardName)
}
