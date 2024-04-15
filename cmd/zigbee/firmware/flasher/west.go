package flasher

import (
	"context"
	"fmt"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/runner"
	"gopkg.in/yaml.v3"
)

type West struct {
	opts map[string]any
}

func (w West) Flash(ctx context.Context, device *config.Device, workDir string) error {
	opts := make([]string, 0, len(w.opts)*2+1)
	opts = append(opts, "flash")

	for opt, val := range w.opts {
		opts = append(opts, "--"+opt, fmt.Sprint(val))
	}

	return runner.NewCmd("west", opts...).Run(
		ctx,
		runner.WithWorkDir(workDir),
		runner.WithToolchainPath(device.General.GetToochainsPath()),
	)
}

func (w *West) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&w.opts)
}
