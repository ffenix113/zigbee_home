package flasher

import (
	"context"

	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/runner"
)

type West struct{}

func (West) Flash(ctx context.Context, device *config.Device, workDir string) error {
	return runner.NewCmd("west", "flash").Run(
		ctx,
		runner.WithWorkDir(workDir),
		runner.WithToolchainPath(device.General.GetToochainsPath()),
	)
}
