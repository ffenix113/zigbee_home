package flasher

import (
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/runner"
)

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
