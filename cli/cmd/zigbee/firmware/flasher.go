package firmware

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/runner"
	"gopkg.in/yaml.v3"
)

var knownFlashers = map[string]any{
	"west":    &West{},
	"nrfutil": NewNRFUtil(),
	"mcuboot": &MCUBoot{},
}

type Flasher interface {
	Flash(ctx context.Context, device *config.Device, workDir string) error
}

func NewFlasher(device *config.Device) Flasher {
	flasherName := "nrfutil"
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

// NRFUtil will flash the board with `nrfutil` application.
// https://docs.zephyrproject.org/latest/boards/arm/nrf52840dongle_nrf52840/doc/index.html#option-1-using-the-built-in-bootloader-only
type NRFUtil struct {
	Port string
}

func NewNRFUtil() *NRFUtil {
	return &NRFUtil{
		Port: "/dev/ttyACM0",
	}
}

func (n NRFUtil) Flash(ctx context.Context, device *config.Device, workDir string) error {
	// 1. Package the app
	packager := runner.NewCmd(
		"nrfutil",
		"pkg",
		"generate",
		"--hw-version",
		"52",
		"--sd-req",
		"0x00",
		"--application",
		"build/zephyr/zephyr.hex",
		"--application-version",
		"1",
		"dfu.zip",
	)
	if err := packager.Run(ctx, runner.WithWorkDir(workDir)); err != nil {
		return fmt.Errorf("create dfu package: %w", err)
	}

	// 2. Flash the app
	flasher := runner.NewCmd(
		"nrfutil",
		"dfu",
		"usb-serial",
		"-pkg",
		"dfu.zip",
		"-p",
		n.Port,
	)
	if err := flasher.Run(ctx, runner.WithWorkDir(workDir)); err != nil {
		return fmt.Errorf("flash the board: %w", err)
	}

	return nil
}

// MCUBoot will flash the board that has mcuboot.
// https://docs.zephyrproject.org/latest/boards/arm/nrf52840dongle_nrf52840/doc/index.html#option-2-using-mcuboot-in-serial-recovery-mode
type MCUBoot struct{}

func (MCUBoot) Flash(ctx context.Context, device *config.Device, workDir string) error {
	opts := []runner.CmdOpt{
		runner.WithWorkDir(workDir),
		runner.WithToolchainPath(getToochainsPath(device)),
	}

	// 1. Sign the firmware
	sign := runner.NewCmd(
		"west",
		"sign",
		"-t",
		"imgtool",
		"--bin",
		"--no-hex",
		"-d",
		"build/zephyr",
		"-B",
		"zephyr.signed.bin",
		"--",
		"--key",
		"some_key.pem",
	)
	if err := sign.Run(ctx, opts...); err != nil {
		return fmt.Errorf("sign app: %w", err)
	}

	// 2. Flash the image
	flash := runner.NewCmd(
		"mcumgr",
		"--conntype",
		"serial",
		"--connstring",
		"dev=/dev/ttyACM0,baud=115200",
		"image",
		"upload",
		"-e",
		"zephyr.signed.bin",
	)
	if err := flash.Run(ctx, opts...); err != nil {
		return fmt.Errorf("flash the board: %w", err)
	}

	// 3. Reset the board
	reset := runner.NewCmd(
		"mcumgr",
		"--conntype",
		"serial",
		"--connstring",
		"dev=/dev/ttyACM0,baud=115200",
		"reset",
	)
	if err := reset.Run(ctx, opts...); err != nil {
		return fmt.Errorf("reset the board: %w", err)
	}

	return nil
}

type West struct{}

func (West) Flash(ctx context.Context, device *config.Device, workDir string) error {
	return runner.NewCmd("west", "flash").Run(
		ctx,
		runner.WithWorkDir(workDir),
		runner.WithToolchainPath(getToochainsPath(device)),
	)
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
