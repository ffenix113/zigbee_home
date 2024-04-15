package config

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"gopkg.in/yaml.v3"
)

type Device struct {
	General General
	Board   Board

	Sensors sensor.Sensors
}

type General struct {
	NCSToolChainBase string `yaml:"ncs_toolchain_base"`
	NCSVersion       string `yaml:"ncs_version"`
	ZephyrBase       string `yaml:"zephyr_base"`

	Manufacturer string
	DeviceName   string
	// Zephyr name for the board
	Board    string
	RunEvery time.Duration
	// ZigbeeChannels will define which endpoints device should try to use.
	// By default device will try all available channels.
	ZigbeeChannels []int `yaml:"zigbee_channels"`
	// Flasher defines the way the board should be flashed.
	Flasher        string
	FlasherOptions map[string]any
}

type Board struct {
	Bootloader *string
	Debug      *extenders.DebugConfig
	IsRouter   bool `yaml:"is_router"`
	LEDs       []types.Pin
	I2C        []extenders.I2CInstance
	UART       []extenders.UARTInstance
}

func ParseFromFile(configPath string) (*Device, error) {
	cfg := &Device{
		General: General{
			RunEvery:         time.Minute,
			NCSToolChainBase: "~/ncs",
			NCSVersion:       "v2.5.0",
		},
	}

	file, err := os.Open(configPath)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	defer file.Close()

	dec := yaml.NewDecoder(file)
	dec.KnownFields(true)

	if err := dec.Decode(cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	// This may contain environment variables,
	// so be kind and try to resolve
	cfg.General.NCSToolChainBase = resolveStringEnv(cfg.General.NCSToolChainBase)
	cfg.General.ZephyrBase = resolveStringEnv(cfg.General.ZephyrBase)

	cfg.PrependCommonClusters()

	return cfg, nil
}

// PrependCommonClusters adds common device clusters as first endpoint.
//
// This allows to have dynamic set of common device clusters,
// such as Identify(server), basic, poll controll, etc.
//
// FIXME: It is mostly a "workaround" to simplify device endpoint generation.
// While the solution is sound to me, the implementation of this function is questionable.
// Should it be here? Should it look like this? Should this common clusters be a sensor,
// rather then converting templates to handle endpoints rather than sensors directly?
func (d *Device) PrependCommonClusters() {
	// Sensors are de-facto our endpoints for now,
	// so prepend common clusters as a sensor.
	d.Sensors = slices.Insert(d.Sensors, 0, sensor.Sensor(base.NewCommonDeviceClusters()))
}

func (g General) GetToochainsPath() (string, string) {
	// If env variables are defined - they have higher priority.
	ncsToolchainPath := os.Getenv("NCS_TOOLCHAIN_BASE")
	ncsVersion := os.Getenv("NCS_VERSION")
	zephyrPath := os.Getenv("ZEPHYR_BASE")

	if ncsVersion == "" {
		ncsVersion = g.NCSVersion
	}

	var locations NCSLocation
	if ncsToolchainPath == "" || zephyrPath == "" {
		var err error
		locations, err = FindNCSLocation(g.NCSToolChainBase, ncsVersion)
		if err != nil {
			log.Panicf("find ncs location: %s", err.Error())
		}

		log.Printf("found toolchain version %q, requested version %q", locations.Version, ncsVersion)
	}

	if ncsToolchainPath == "" {
		ncsToolchainPath = locations.NCS
	}

	if zephyrPath == "" {
		zephyrPath = locations.Zephyr
	}

	return ncsToolchainPath, zephyrPath
}

func resolveStringEnv(input string) string {
	if strings.HasPrefix(input, "~/") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("could not resolve user home dir: %s", err.Error()))
		}

		input = strings.Replace(input, "~/", userHome+"/", 1)
	}

	return os.ExpandEnv(input)
}
