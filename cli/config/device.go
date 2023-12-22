package config

import (
	"github.com/ffenix113/zigbee_home/cli/sensor"
	"github.com/ffenix113/zigbee_home/cli/types"
)

type Device struct {
	General General
	Board   Board

	Sensors sensor.Sensors
}

type General struct {
	Manufacturer string
	DeviceName   string
	// Zephyr name for the board
	Board string
	// Flasher defines the way the board should be flashed.
	Flasher        string
	FlasherOptions map[string]any
}

type Board struct {
	Uart Uart
}

type Uart struct {
	Rx, Tx types.Pin
}
