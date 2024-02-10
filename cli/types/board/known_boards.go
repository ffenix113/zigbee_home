package board

import (
	"slices"

	"github.com/ffenix113/zigbee_home/cli/types/appconfig"
)

type Bootloader struct {
	PM     *PMConfig
	Config []appconfig.ConfigValue
}

type Section struct {
	Address int
	Size    int
}

type PMConfig struct {
	Flash Section
	SRAM  Section
}

func (s Section) EndAddress() int {
	return s.Address + s.Size
}

func BoardBootloaderConfig(board string) (*Bootloader, string) {
	var boardBootloader string
	for bootloader, boards := range bootloaderBoards {
		if _, found := slices.BinarySearch(boards, board); found {
			boardBootloader = bootloader
			break
		}
	}

	return BootloaderConfig(boardBootloader), boardBootloader
}

func BootloaderConfig(bootloader string) *Bootloader {
	// By default will return nil, if `bootloader` is empty.
	// This would be enough to not include any bootloader configuration
	return knownBootloaders[bootloader]
}

var bootloaderBoards = func() map[string][]string {
	// Values in this map can be added in any order,
	// as they will be ordered on the next step.
	bootloadersMap := map[string][]string{
		"arduino": {
			"arduino_nano_33_ble",
		},
		"nrf52_legacy": {
			"nrf52840dongle_nrf52840",
		},
		"adafruit_nrf52_sd132":    {},
		"adafruit_nrf52_sd140_v6": {},
		"adafruit_nrf52_sd140_v7": {},
	}

	for _, boards := range bootloadersMap {
		slices.Sort(boards)
	}

	return bootloadersMap
}()

var knownBootloaders = map[string]*Bootloader{
	"nrf52_legacy": { // Legacy NRF52 bootloader, provided with nRF5 SDK. Also used in nRF52840 Dongle till now for some reason.
		PM: &PMConfig{
			Flash: Section{
				Address: 0xe0000,
				Size:    0x20000,
			},
			SRAM: Section{
				Address: 0x20000000,
				Size:    0x400,
			},
		},
	},
	"arduino": { // Bossac bootloader
		PM: &PMConfig{
			Flash: Section{
				Address: 0x0,
				Size:    0x10000,
			},
			// App is launching and starting to work without SRAM config,
			// so we will not define it(for now at least).
			// Maybe there will be issues in the future, and they will be fixed here.
		},
	},
	"adafruit_nrf52_sd132": {
		PM: &PMConfig{
			// https://infocenter.nordicsemi.com/topic/sdk_nrf5_v17.1.0/lib_bootloader.html, "Memory layout" section
			Flash: Section{
				Address: 0x0,
				Size:    0x26000,
			},
		},
		Config: []appconfig.ConfigValue{
			appconfig.NewValue("CONFIG_BUILD_OUTPUT_UF2").Required(appconfig.Yes),
			appconfig.NewValue("CONFIG_BOOTLOADER_BOSSA").Required(appconfig.No),
			appconfig.NewValue("CONFIG_FLASH_LOAD_OFFSET").Required("0x26000"),
		},
	},
	"adafruit_nrf52_sd140_v6": {
		PM: &PMConfig{
			// https://infocenter.nordicsemi.com/topic/sdk_nrf5_v17.1.0/lib_bootloader.html, "Memory layout" section
			Flash: Section{
				Address: 0x0,
				Size:    0x26000,
			},
		},
		Config: []appconfig.ConfigValue{
			appconfig.NewValue("CONFIG_BUILD_OUTPUT_UF2").Required(appconfig.Yes),
			// This option is to prevent use of default flash load offset.
			// May be specific to some boards, but we need it just to be safe.
			appconfig.NewValue("CONFIG_BOOTLOADER_BOSSA").Required(appconfig.No),
			appconfig.NewValue("CONFIG_FLASH_LOAD_OFFSET").Required("0x26000"),
		},
	},
	"adafruit_nrf52_sd140_v7": {
		PM: &PMConfig{
			// https://infocenter.nordicsemi.com/topic/sdk_nrf5_v17.1.0/lib_bootloader.html, "Memory layout" section
			Flash: Section{
				Address: 0x0,
				Size:    0x27000,
			},
		},
		Config: []appconfig.ConfigValue{
			appconfig.NewValue("CONFIG_BUILD_OUTPUT_UF2").Required(appconfig.Yes),
			appconfig.NewValue("CONFIG_BOOTLOADER_BOSSA").Required(appconfig.No),
			appconfig.NewValue("CONFIG_FLASH_LOAD_OFFSET").Required("0x27000"),
		},
	},
}
