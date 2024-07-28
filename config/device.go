package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/ffenix113/zigbee_home/types/yamlstrict"
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

	Manufacturer string `yaml:"manufacturer"`
	DeviceName   string `yaml:"device_name"`
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
	Bootloader         *string
	Debug              *extenders.DebugConfig
	IsRouter           bool   `yaml:"is_router"`
	FactoryResetButton string `yaml:"factory_reset_button"`
	NetworkStateLED    string `yaml:"network_state_led"`
	LEDs               types.PinWithIDSlice
	// Buttons provide definitions(or references) to board buttons.
	// They will be used in other configuration places to
	// reference specific button.
	Buttons types.PinWithIDSlice
	I2C     []extenders.I2CInstance
	UART    []extenders.UARTInstance
}

func ParseFromFile(configPath string) (*Device, error) {
	cfg := &Device{
		General: General{
			RunEvery: time.Minute,
			NCSToolChainBase: func() string {
				if runtime.GOOS == "windows" {
					return "C:\\ncs"
				}
				return "~/ncs"
			}(),
			NCSVersion:   "v2.6.1",
			Manufacturer: "FFexix113",
			DeviceName:   "dongle",
		},
	}

	file, err := os.Open(configPath)
	if err != nil {
		return cfg, fmt.Errorf("open config file: %w", err)
	}

	defer file.Close()

	cfg, err = ParseFromReader(cfg, file)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	return cfg, nil
}

func ParseFromReader(defConfig *Device, rdr io.Reader) (*Device, error) {
	dec := yaml.NewDecoder(rdr)
	dec.KnownFields(false)

	if err := dec.Decode(defConfig); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// This may contain environment variables,
	// so be kind and try to resolve
	defConfig.General.NCSToolChainBase = resolveStringEnv(defConfig.General.NCSToolChainBase)

	defConfig.General.ZephyrBase = resolveStringEnv(defConfig.General.ZephyrBase)

	defConfig.PrependCommonClusters()

	return defConfig, nil
}

// UnamrshalYAML is implemented to intercept the original
// configuration file and resolve any known tags inside.
func (d *Device) UnmarshalYAML(node *yaml.Node) error {
	resolver := newTagsResolver()
	if err := resolver.resolve(node, 1); err != nil {
		return fmt.Errorf("resolve tags: %w", err)
	}

	type dev Device

	if err := yamlstrict.Unmarshal((*dev)(d), node); err != nil {
		return fmt.Errorf("unamrshal config: %w", err)
	}

	return nil
}

// PrependCommonClusters adds common device clusters as first endpoint.
//
// This allows to have dynamic set of common device clusters,
// such as Identify(server), basic, poll control, etc.
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

func (g General) GetToochainsPath() NCSLocation {
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
			log.Fatalf("find ncs location: %s", err.Error())
		}

		log.Printf("found toolchain version %q, requested version %q", locations.Version, ncsVersion)
	}

	if ncsToolchainPath == "" {
		ncsToolchainPath = locations.NCS
	}

	if zephyrPath == "" {
		zephyrPath = locations.Zephyr
	}

	return NCSLocation{
		Version: locations.Version,
		NCS:     ncsToolchainPath,
		Zephyr:  zephyrPath,
	}
}

func resolveStringEnv(input string) string {
	if strings.HasPrefix(input, "~/") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			panic("could not resolve user home dir: " + err.Error())
		}

		input = strings.Replace(input, "~/", userHome+"/", 1)
	}

	return os.ExpandEnv(input)
}
