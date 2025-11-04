package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/ffenix113/zigbee_home/sensor/base"
	"github.com/ffenix113/zigbee_home/templates/extenders"
	"github.com/ffenix113/zigbee_home/types"
	"github.com/ffenix113/zigbee_home/types/board"
	"github.com/ffenix113/zigbee_home/types/sensor"
	"github.com/ffenix113/zigbee_home/types/yamlstrict"
	"gopkg.in/yaml.v3"
)

type Device struct {
	// This field is necessary for import resolution step,
	// if it would be possible to provide it in any other way -
	// this field should be removed.
	// An ugly hack, so to speak.
	configPath string

	General General
	Board   Board

	Sensors sensor.Sensors
}

type General struct {
	NCS NCS
	// NCSToolChainBase string `yaml:"ncs_toolchain_base"`
	// NCSVersion       string `yaml:"ncs_version"`
	// ZephyrBase       string `yaml:"zephyr_base"`

	// TemplatesPath allows to override default templates path.
	// This would allow to use custom templates while developing, for example.
	TemplatesPath string `yaml:"templates_path"`

	Manufacturer string `yaml:"manufacturer"`
	DeviceName   string `yaml:"device_name"`
	// Zephyr name for the board
	Board string
	SoC   string

	RunEvery time.Duration
	// ZigbeeChannels will define which endpoints device should try to use.
	// By default device will try all available channels.
	ZigbeeChannels []int `yaml:"zigbee_channels"`
	// TrustCenterKey is an optional configuration that will allow to
	// set trust center key, which in turn would allow connecting to
	// specific Zigbee hubs, like Philips Hue Bridge.
	//
	// This configuration option requires device to be set up
	// as router or coordinator. End device does not support this option.
	//
	// Note: This option is EXPERIMENTAL
	TrustCenterKey string `yaml:"trust_center_key"`
	// Flasher defines the way the board should be flashed.
	Flasher        string
	FlasherOptions map[string]any
}

type Board struct {
	Bootloader *string
	Debug      *extenders.DebugConfig
	// EnableWatchdog will setup watchdog timer that
	// will reset SoC if it was not "fed" (i.e. pinged)
	// in some time.
	//
	// This option is set to true by default,
	// and can only be disabled.
	//
	// It is always disabled if configuration is set in debug mode.
	EnableWatchdog     bool   `yaml:"enable_watchdog"`
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

type NCS struct {
	ToolchainBasePath string `yaml:"toolchain_base_path"`
	// ToolchainVersion  string `yaml:"toolchain_version"`

	SDKBasePath string `yaml:"sdk_base_path"`
	SDKVersion  string `yaml:"sdk_version"`
}

func ParseFromFile(configPath string) (*Device, error) {
	var minimumSDKVersion = types.NewSemver(2, 6, 0)

	cfg := &Device{
		configPath: filepath.Dir(configPath),
		General: General{
			RunEvery: time.Minute,
			NCS: NCS{
				SDKVersion: minimumSDKVersion.String(),
				// ToolchainVersion: minimumNCSVersion.String(),
				// Toolchain path will be resolved based on SDK base path,
				// if it is not specifically provided.
				SDKBasePath: func() string {
					switch runtime.GOOS {
					case "windows":
						return "C:\\ncs"
					case "darwin":
						return "/opt/nordic/ncs"
					default:
						return os.ExpandEnv("$HOME/ncs")
					}
				}(),
			},
			Manufacturer: "zigbee_home",
			DeviceName:   "device",
		},
		Board: Board{
			EnableWatchdog: true,
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

	if cfg.General.SoC == "" {
		cfg.General.SoC, err = ResolveBoardSoC(cfg)
		if err != nil {
			return nil, fmt.Errorf("resolve board soc: %w", err)
		}
	}

	if err := ValidateConfiguration(cfg); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	return cfg, nil
}

func ParseFromReader(defConfig *Device, rdr io.Reader) (*Device, error) {
	dec := yaml.NewDecoder(rdr)
	dec.KnownFields(true)

	if err := dec.Decode(defConfig); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	defConfig.PrependCommonClusters()

	return defConfig, nil
}

func ResolveBoardSoC(conf *Device) (string, error) {
	boardNameParts := strings.Split(conf.General.Board, "/")

	for _, namePart := range boardNameParts {
		if board.IsKnownSoC(namePart) {
			return namePart, nil
		}
	}

	// Resolve board soc through Zephyr board index(hopefully board will have only one supported soc):
	// boardDir := cmd.Exec(`west boards --board {boardNameParts[0]} -f "{dir}"`)
	// boardDefinitionPath := boardDir + "/board.yml"
	// boardDefinition := yaml.Parse(boardDefinitionPart)
	// for _, socDef := range boardDefition.["board"]["socs"] {
	// 	if board.IsKnownSoC(socDef["name"]) {
	// 		return soc, nil
	// 	}
	// }

	return "", nil
}

// UnamrshalYAML is implemented to intercept the original
// configuration file and resolve any known tags inside.
func (d *Device) UnmarshalYAML(node *yaml.Node) error {
	resolver := newTagsResolver(d.configPath)
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
	// If we already have env set - don't do anything.
	if DoNotSetupEnv() {
		return NCSLocation{}
	}

	// If env variables are defined - they have higher priority.
	ncsToolchainPath := os.Getenv("NCS_TOOLCHAIN_BASE")
	ncsSDKVersion := os.Getenv("NCS_VERSION")
	zephyrSDKPath := os.Getenv("ZEPHYR_BASE")

	if ncsSDKVersion == "" {
		ncsSDKVersion = g.NCS.SDKVersion
	}

	var locations NCSLocation

	if ncsToolchainPath == "" || zephyrSDKPath == "" {
		var err error
		locations, err = FindNCSLocation(g.NCS.SDKBasePath, ncsSDKVersion)

		if err != nil {
			log.Fatalf("find ncs location: %s", err.Error())
		}

		log.Printf("found sdk version %q, toolchain version %q, requested version %q", locations.SDKVersion, locations.ToolchainVersion, ncsSDKVersion)
	}

	if ncsToolchainPath == "" {
		ncsToolchainPath = locations.ToolchainPath
	}

	if zephyrSDKPath == "" {
		zephyrSDKPath = locations.SDKPath
	}

	return NCSLocation{
		SDKVersion:       locations.SDKVersion,
		ToolchainVersion: locations.ToolchainVersion,
		ToolchainPath:    ncsToolchainPath,
		SDKPath:          zephyrSDKPath,
	}
}

// ValidateConfiuration checks device configuration as much as it can
// to provide meaningful information about errors in configuration.
func ValidateConfiguration(cfg *Device) error {
	if cfg.General.TrustCenterKey != "" {
		log.Println("EXPERIMENTAL: Trust center key configuration is experimental and may not work as expected")

		// Verify that key is of correct length.
		// 32 chars + 15 colons
		const keyLength = 32 + 15
		if providedKeyLen := len(cfg.General.TrustCenterKey); providedKeyLen != keyLength {
			return fmt.Errorf("trust center key has wrong length: want %d, have %d", keyLength, providedKeyLen)
		}

		if colonsCount := strings.Count(cfg.General.TrustCenterKey, ":"); colonsCount != 15 {
			return fmt.Errorf("trust center key bytes shold be separated by colon: has %d colons, want 15", colonsCount)
		}

		// Trust center key is supported only for router & coordinator roles.
		if !cfg.Board.IsRouter {
			return errors.New("trust center key is provided, but board is not configured to be router")
		}
	}

	return nil
}

func DoNotSetupEnv() bool {
	_, ok := os.LookupEnv("NO_SETUP_ENV")

	return ok
}
