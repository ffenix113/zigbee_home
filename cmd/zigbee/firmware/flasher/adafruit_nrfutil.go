package flasher

import (
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/runner"
)

// AdafruitNRFUtil will flash the board with `nrfutil` application.
// For now it only works if `adafruit-nrfutil` application is in current PATH.
//
// So, by default, if it is installed into venv - it will not work.
//
// https://docs.zephyrproject.org/latest/boards/arm/nrf52840dongle_nrf52840/doc/index.html#option-1-using-the-built-in-bootloader-only
type AdafruitNRFUtil struct {
	Port string
}

func NewAdafruitNRFUtil() *AdafruitNRFUtil {
	return &AdafruitNRFUtil{
		Port: "/dev/ttyACM0",
	}
}

func (n AdafruitNRFUtil) Flash(ctx context.Context, device *config.Device, workDir string) error {
	// 1. Package the app
	packager := runner.NewCmd(
		"adafruit-nrfutil",
		"dfu",
		"genpkg",
		"--dev-type",
		"0x0052",
		"--application",
		"build/zephyr/zephyr.hex",
		"dfu.zip",
	)
	if err := packager.Run(ctx, runner.WithWorkDir(workDir)); err != nil {
		return fmt.Errorf("create dfu package: %w", err)
	}

	// 2. Flash the app
	flasher := runner.NewCmd(
		"adafruit-nrfutil",
		"dfu",
		"serial",
		"--package",
		"dfu.zip",
		"-p",
		n.Port,
		"--singlebank",
		"-b",
		"115200",
	)
	if err := flasher.Run(ctx, runner.WithWorkDir(workDir)); err != nil {
		return fmt.Errorf("flash the board: %w", err)
	}

	return nil
}
