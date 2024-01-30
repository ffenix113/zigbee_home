package board

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

func BoardPMConfig(board string) *PMConfig {
	return knownBoardConfigs[board]
}

var knownBoardConfigs = map[string]*PMConfig{
	"nrf52840dongle_nrf52840": {
		Flash: Section{
			Address: 0xe0000,
			Size:    0x20000,
		},
		SRAM: Section{
			Address: 0x20000000,
			Size:    0x400,
		},
	},
	"arduino_nano_33_ble": {
		Flash: Section{
			Address: 0x0,
			Size:    0x10000,
		},
		// App is launching and starting to work without SRAM config,
		// so we will not define it(for now at least).
		// Maybe there will be issues in the future, and they will be fixed here.
	},
}
