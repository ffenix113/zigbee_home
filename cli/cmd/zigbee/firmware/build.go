package firmware

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/generate"
	"github.com/ffenix113/zigbee_home/cli/runner"
	"github.com/urfave/cli/v2"
)

func buildCmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build the firmware",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "only-generate",
			},
		},
		Action: buildFirmware,
	}
}

func buildFirmware(ctx *cli.Context) error {
	cfg, err := parseConfig(ctx)
	if err != nil {
		return fmt.Errorf("prepare config: %w", err)
	}

	if cfg.General.Board == "" {
		return fmt.Errorf("board name cannot be empty")
	}

	// Will work in the future.
	workDir, err := filepath.Abs(ctx.String("workdir"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if workDir == "" {
		workDir = "."
	}

	generator, err := generate.NewGenerator(cfg)
	if err != nil {
		return fmt.Errorf("new generator: %w", err)
	}

	if err := generator.Generate(workDir, cfg); err != nil {
		return fmt.Errorf("generate base: %w", err)
	}

	if !ctx.Bool("only-generate") {
		return runBuild(ctx.Context, cfg, workDir)
	}

	return nil
}

func parseConfig(ctx *cli.Context) (*config.Device, error) {
	configPath := ctx.String("config")
	if configPath == "" {
		return nil, errors.New("config path cannot be empty (it is set by default)")
	}

	conf, err := config.ParseFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return conf, nil
}

func runBuild(ctx context.Context, device *config.Device, workDir string) error {
	build := runner.NewCmd(
		"west",
		"build",
		"--pristine", // For now let's always build Pristine.
		"--board", device.General.Board,
		"--no-sysbuild", // https://docs.zephyrproject.org/latest/build/sysbuild/index.html
		"--build-dir", workDir+"/build",
		workDir,
		"--",
		"-DNCS_TOOLCHAIN_VERSION=NONE",
		fmt.Sprintf("-DCONF_FILE=%s/prj.conf", workDir),
		fmt.Sprintf("-DDTC_OVERLAY_FILE=%s/app.overlay", workDir),
	)

	if err := build.Run(ctx, runner.WithToolchainPath(device.General.GetToochainsPath())); err != nil {
		return fmt.Errorf("build firmware: %w", err)
	}

	return nil
}
