package firmware

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func flashCmd() *cli.Command {
	return &cli.Command{
		Name:  "flash",
		Usage: "flash the firmware",
		Flags: configFileFlags(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configFileName := getConfigFile(cmd)

			cfg, err := parseConfig(configFileName, cmd.StringSlice("override"))
			if err != nil {
				return fmt.Errorf("prepare config: %w", err)
			}

			// Will work in the future.
			workDir := cmd.String("workdir")
			if workDir == "" {
				workDir = "."
			}

			// generator := generate.NewGenerator(cfg)

			// if err := generator.Generate(workDir); err != nil {
			// 	return fmt.Errorf("generate base: %w", err)
			// }

			flasher := NewFlasher(cfg)

			return flasher.Flash(ctx, cfg, workDir)
		},
	}
}
