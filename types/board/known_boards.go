package board

import (
	"slices"
	"strconv"

	"github.com/ffenix113/zigbee_home/types/appconfig"
)

type Bootloader struct {
	PM     []Section
	Config []appconfig.ConfigValue
}

const (
	RegionFlashPrimary = "flash_primary"
	RegionSramPrimary  = "sram_primary"
)

type Section struct {
	Name    string
	Address int
	Size    int
	Region  string
}

func (s Section) EndAddress() int {
	return s.Address + s.Size
}

func BootloaderConfigFromBoard(board string) (*Bootloader, string) {
	var boardBootloader string

	for bootloader, boards := range bootloaderBoards() {
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
	return knownBootloaders()[bootloader]
}

func bootloaderBoards() map[string][]string {
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
		"adafruit_nrf52_sd140_v7": {
			"xiao_ble",
		},
	}

	for _, boards := range bootloadersMap {
		slices.Sort(boards)
	}

	return bootloadersMap
}

func knownBootloaders() map[string]*Bootloader {
	adafruitConfig := func(appStartAddr int) *Bootloader {
		appStartAddrHex := "0x" + strconv.FormatInt(int64(appStartAddr), 16)

		return &Bootloader{
			PM: []Section{
				{
					Name:    "empty_app_offset",
					Address: 0x0,
					Size:    appStartAddr,
					Region:  RegionFlashPrimary,
				},
				{
					// This section is needed because without it
					// the app will be from appStartAddr till 0xf7000,
					// which is right in the middle of bootloader.
					// Next section will be for zboss and it will start
					// at 0xf7000, which will overwrite part of bootloader
					// resulting in a dead board after power-cycle,
					// but can be resurected if bootloader is re-flashed.
					Name:    "empty_after_zboss_offset",
					Address: 0xf4000,
					Size:    0xc000,
					Region:  RegionFlashPrimary,
				},
			},
			Config: []appconfig.ConfigValue{
				appconfig.NewValue("CONFIG_BUILD_OUTPUT_UF2").Required(appconfig.Yes),
				// This option is to prevent use of default flash load offset.
				// May be specific to some boards, but we need it just to be safe.
				appconfig.NewValue("CONFIG_BOOTLOADER_BOSSA").Required(appconfig.No),
				appconfig.NewValue("CONFIG_FLASH_LOAD_OFFSET").Required(appStartAddrHex),
			},
		}
	}

	return map[string]*Bootloader{
		// Legacy NRF52 bootloader, provided with nRF5 SDK. Also used in nRF52840 Dongle till now for some reason.
		"nrf52_legacy": {
			// This configuration is based on
			// https://github.com/martelmy/NCS_examples/blob/3cd874c68e13565ca89a082c67af54f48b704192/zigbee/light_bulb_dongle/pm_static_nrf52840dongle_nrf52840.yml
			// , based on the answer here:
			// https://devzone.nordicsemi.com/f/nordic-q-a/87169/zigbee-example-for-nrf52840-dongle/364086
			//
			// Thank you Marte Myrvold(https://github.com/martelmy)
			PM: []Section{
				{
					Name:    "empty_bootloader",
					Address: 0xe0000,
					Size:    0x20000,
					Region:  RegionFlashPrimary,
				},
				{
					Name:    "empty_sram_bootloader",
					Address: 0x20000000,
					Size:    0x400,
					Region:  RegionSramPrimary,
				},
			},
		},
		// Bossac bootloader
		"arduino": {
			PM: []Section{
				{
					Name:    "empty_before_app",
					Address: 0x0,
					Size:    0x10000,
					Region:  RegionFlashPrimary,
				},
				// App is launching and starting to work without SRAM config,
				// so we will not define it(for now at least).
				// Maybe there will be issues in the future, and they will be fixed here.
			},
		},
		"adafruit_nrf52_sd132":    adafruitConfig(0x26000),
		"adafruit_nrf52_sd140_v6": adafruitConfig(0x26000),
		"adafruit_nrf52_sd140_v7": adafruitConfig(0x27000),
	}
}
