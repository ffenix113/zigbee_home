package config

import (
	"time"

	"github.com/ffenix113/zigbee_home/cli/templates/extenders"
	"github.com/ffenix113/zigbee_home/cli/types"
	"github.com/ffenix113/zigbee_home/cli/types/sensor"
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
	Board    string
	RunEvery time.Duration
	// Flasher defines the way the board should be flashed.
	Flasher        string
	FlasherOptions map[string]any
}

type Board struct {
	DebugLog bool
	I2C      *extenders.I2C
}

type Uart struct {
	Rx, Tx types.Pin
}
