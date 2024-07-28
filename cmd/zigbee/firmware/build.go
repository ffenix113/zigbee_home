package firmware

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/ffenix113/zigbee_home/generate"
	"github.com/ffenix113/zigbee_home/runner"
	"github.com/urfave/cli/v2"
)

const filenameArg = "config"

// BuildConfig is aimed to provide build information
// independently of how that information was obtained.
type BuildConfig struct {
	WorkDir      string
	ConfigFile   string
	OnlyGenerate bool
	ClearWorkDir bool
}

func buildCmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build the firmware",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "only-generate",
			},
			&cli.BoolFlag{
				Name: "clear-work-dir",
			},
		},
		Action: func(ctx *cli.Context) error {
			buildCtx, err := newBuildConfigFromCLI(ctx)
			if err != nil {
				return fmt.Errorf("build context: %w", err)
			}

			return BuildFirmware(ctx.Context, buildCtx)
		},
	}
}

func newBuildConfigFromCLI(ctx *cli.Context) (BuildConfig, error) {
	workdir, err := getWorkdir(ctx)
	if err != nil {
		return BuildConfig{}, fmt.Errorf("get workdir: %w", err)
	}

	return BuildConfig{
		WorkDir:      workdir,
		ConfigFile:   getConfigFile(ctx),
		OnlyGenerate: ctx.Bool("only-generate"),
		ClearWorkDir: ctx.Bool("clear-work-dir"),
	}, nil
}

func BuildFirmware(ctx context.Context, buildConfig BuildConfig) error {
	cfg, err := parseConfig(buildConfig.ConfigFile)
	if err != nil {
		return fmt.Errorf("prepare config: %w", err)
	}

	if cfg.General.Board == "" {
		return fmt.Errorf("board name cannot be empty")
	}

	if err := GenerateFirmwareFiles(ctx, buildConfig.WorkDir, buildConfig.ClearWorkDir, cfg); err != nil {
		return fmt.Errorf("generate firmware files: %w", err)
	}

	if !buildConfig.OnlyGenerate {
		return runBuild(ctx, cfg, buildConfig.WorkDir)
	}

	return nil
}

func GenerateFirmwareFiles(ctx context.Context, workDir string, shouldClearWorkDir bool, cfg *config.Device) error {
	generator, err := generate.NewGenerator(cfg)
	if err != nil {
		return fmt.Errorf("new generator: %w", err)
	}

	if shouldClearWorkDir {
		if err := clearWorkDir(workDir); err != nil {
			return err
		}
	}

	if err := generator.Generate(workDir, cfg); err != nil {
		return fmt.Errorf("generate base: %w", err)
	}

	return nil
}

func parseConfig(configPath string) (*config.Device, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty (it is set by default)")
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("make config path %q absolute: %w", configPath, err)
	}

	conf, err := config.ParseFromFile(absConfigPath)
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

	toolchainsPath := device.General.GetToochainsPath()
	if err := build.Run(ctx, runner.WithToolchainPath(toolchainsPath.NCS, toolchainsPath.Zephyr)); err != nil {
		return fmt.Errorf("build firmware: %w", err)
	}

	return nil
}

func clearWorkDir(workDir string) error {
	return filepath.WalkDir(workDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			var pathErr *os.PathError
			if !errors.As(err, &pathErr) {
				return err
			}

			if errors.Is(pathErr, os.ErrNotExist) {
				return nil
			}

			return err
		}

		if path == workDir {
			return nil
		}

		return os.RemoveAll(path)
	})
}

func getWorkdir(ctx *cli.Context) (string, error) {
	workDir, err := filepath.Abs(ctx.String("workdir"))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if workDir == "" {
		workDir = "."
	}

	// This will make sure that workdir uses slashes as path separators even on windows,
	// which will be fine for cmake.
	workDir = filepath.ToSlash(workDir)

	return workDir, nil
}

func getConfigFile(ctx *cli.Context) string {
	if ctx.IsSet(filenameArg) {
		return ctx.String(filenameArg)
	}

	preferences := []string{"zigbee.yaml", "zigbee.yml"}
	for _, preference := range preferences {
		if _, err := os.Stat(preference); err == nil {
			if preference == "zigbee.yml" {
				log.Println("Default config file name changed to 'zigbee.yaml', please change name of your configuration file.")
			}

			return preference
		}
	}

	// If both files don't exist - return default value.
	return "zigbee.yaml"
}
