package flasher

import (
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/runner"
)

// MCUBoot will flash the board that has mcuboot.
// https://docs.zephyrproject.org/latest/boards/arm/nrf52840dongle_nrf52840/doc/index.html#option-2-using-mcuboot-in-serial-recovery-mode
type MCUBoot struct{}

func (MCUBoot) Flash(ctx context.Context, device *config.Device, workDir string) error {
	toolchainsPath := device.General.GetToochainsPath()
	opts := []runner.CmdOpt{
		runner.WithWorkDir(workDir),
		runner.WithToolchainPath(toolchainsPath.NCS, toolchainsPath.Zephyr),
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
