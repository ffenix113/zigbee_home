package firmware

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ffenix113/zigbee_home/cli/config"
	"github.com/ffenix113/zigbee_home/cli/generate"
	"github.com/ffenix113/zigbee_home/cli/runner"
	"github.com/urfave/cli/v2"

	"gopkg.in/yaml.v3"
)

func buildCmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "build the firmware",
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
	workDir := ctx.String("workdir")
	if workDir == "" {
		workDir = "."
	}

	generator := generate.NewGenerator(cfg)

	if err := generator.Generate(workDir); err != nil {
		return fmt.Errorf("generate base: %w", err)
	}

	return runBuild(ctx.Context, cfg, workDir)
}

func parseConfig(ctx *cli.Context) (config.Device, error) {
	configPath := ctx.String("config")
	if configPath == "" {
		return config.Device{}, errors.New("config path cannot be empty (it is set by default)")
	}

	conf, err := parseConfigFile(configPath)
	if err != nil {
		return config.Device{}, fmt.Errorf("parse config file: %w", err)
	}

	return conf, nil
}

func parseConfigFile(configPath string) (config.Device, error) {
	var cfg config.Device

	file, err := os.Open(configPath)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	defer file.Close()

	dec := yaml.NewDecoder(file)
	dec.KnownFields(true)

	if err := dec.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func runBuild(ctx context.Context, device config.Device, workDir string) error {
	build := runner.NewCmd(
		"west",
		"build",
		"-p", // For now let's always build Pristine.
		"-b",
		device.General.Board,
	)

	if err := build.Run(ctx, workDir); err != nil {
		return fmt.Errorf("build firmware: %w", err)
	}

	return nil
}
