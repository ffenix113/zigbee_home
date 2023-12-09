package firmware

import (
	"fmt"

	"github.com/ffenix113/zigbee_home/cli/generate"
	"github.com/urfave/cli/v2"
)

func flashCmd() *cli.Command {
	return &cli.Command{
		Name:  "flash",
		Usage: "flash the firmware",
		Action: func(ctx *cli.Context) error {
			cfg, err := parseConfig(ctx)
			if err != nil {
				return fmt.Errorf("prepare config: %w", err)
			}

			// Will work in the future.
			workDir := ctx.String("workdir")
			if workDir == "" {
				workDir = "."
			}

			generator := generate.NewGenerator(cfg)

			if err := generator.Generate(workDir); err != nil {
				return fmt.Errorf("generate base: %w", err)
			}

			flasher := NewFlasher(cfg.General)

			return flasher.Flash(ctx.Context, cfg.General, workDir)
		},
	}
}
