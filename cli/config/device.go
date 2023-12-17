package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ffenix113/zigbee_home/cli/sensor"
	"gopkg.in/yaml.v3"
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
	Rx, Tx Pin
}

type Pin struct {
	Port int
	Pin  int
}

var pinRegex = regexp.MustCompile(`^[01]\.[0-9][0-9]$`)

func (p *Pin) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("pin type should be scalar, but is %q", value.LongTag())
	}

	if !pinRegex.MatchString(value.Value) {
		return fmt.Errorf("pin definition must be in a form of X.XX, where X is a number")
	}

	parts := strings.Split(value.Value, ".")

	p.Port, _ = strconv.Atoi(parts[0])
	p.Pin, _ = strconv.Atoi(parts[1])

	return nil
}
