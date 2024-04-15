package firmware

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/cmd/zigbee/firmware/flasher"
	"github.com/ffenix113/zigbee_home/config"
	"gopkg.in/yaml.v3"
)

var knownFlashers = map[string]any{
	"adafruit": flasher.NewAdafruitNRFUtil(),
	"west":     &flasher.West{},
	"nrfutil":  flasher.NewNRFUtil(),
	"mcuboot":  &flasher.MCUBoot{},
}

type Flasher interface {
	Flash(ctx context.Context, device *config.Device, workDir string) error
}

func NewFlasher(device *config.Device) Flasher {
	// Sane default to flash with `west`.
	flasherName := "west"
	if device.General.Flasher != "" {
		flasherName = device.General.Flasher
	}

	flasher, ok := knownFlashers[flasherName]
	if !ok {
		panic(fmt.Sprintf("flasher %q is not defined", flasherName))
	}

	if len(device.General.FlasherOptions) != 0 {
		if err := readFlasherConfig(device.General.FlasherOptions, flasher); err != nil {
			panic(fmt.Sprintf("read flasher config: %s", err.Error()))
		}
	}

	return flasher.(Flasher)
}

func readFlasherConfig(opts map[string]any, flasher any) error {
	bts, err := yaml.Marshal(opts)
	if err != nil {
		return fmt.Errorf("marshal flasher config: %w", err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(bts))
	dec.KnownFields(true)

	if err := dec.Decode(flasher); err != nil {
		return fmt.Errorf("decode flasher options: %w", err)
	}

	return nil
}
