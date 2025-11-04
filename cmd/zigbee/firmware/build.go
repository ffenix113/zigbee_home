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
	"github.com/urfave/cli/v3"
)

const filenameArg = "config"

// BuildConfig is aimed to provide build information
// independently of how that information was obtained.
type BuildConfig struct {
	WorkDir        string
	ConfigFilePath string
	OnlyGenerate   bool
	ClearWorkDir   bool
}

func buildCmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build the firmware",
		Flags: append(configFileFlags(),
			&cli.BoolFlag{
				Name:  "only-generate",
				Usage: "Only generate source files, do not build them.",
			},
			&cli.BoolFlag{
				Name:  "clear-work-dir",
				Usage: "Will remove all files from specified workspace. Be sure to select correct workspace before running this command!",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			buildCtx, err := newBuildConfigFromCLI(cmd)
			if err != nil {
				return fmt.Errorf("build context: %w", err)
			}

			return BuildFirmware(ctx, buildCtx, cmd.StringSlice("override"))
		},
	}
}

func newBuildConfigFromCLI(cmd *cli.Command) (BuildConfig, error) {
	workdir, err := getWorkdir(cmd)
	if err != nil {
		return BuildConfig{}, fmt.Errorf("get workdir: %w", err)
	}

	return BuildConfig{
		WorkDir:        workdir,
		ConfigFilePath: getConfigFile(cmd),
		OnlyGenerate:   cmd.Bool("only-generate"),
		ClearWorkDir:   cmd.Bool("clear-work-dir"),
	}, nil
}

func BuildFirmware(ctx context.Context, buildConfig BuildConfig, overrides []string) error {
	cfg, err := parseConfig(buildConfig.ConfigFilePath, overrides)
	if err != nil {
		return fmt.Errorf("prepare config: %w", err)
	}

	if cfg.General.Board == "" {
		return fmt.Errorf("board name cannot be empty")
	}

	if _, err = os.Stat(buildConfig.WorkDir); err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			if err := os.Mkdir(buildConfig.WorkDir, os.ModeDir|0o755); err != nil {
				return fmt.Errorf("create workdir: %w", err)
			}
		case errors.Is(err, os.ErrPermission):
			return fmt.Errorf("access workdir permission error: %w", err)
		default:
			return fmt.Errorf("stat workdir: %w", err)
		}
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

func runBuild(ctx context.Context, device *config.Device, workDir string) error {
	build := runner.NewCmd(
		"west",
		"build",
		"--pristine", // For now let's always build Pristine.
		"--board", device.General.Board,
		// Enable sysbuild, as it is required for newer nRF Connect SDK
		// and it will allow to build MCUBoot image as well.
		// It is also used to build network image for nRF53 series.
		"--sysbuild", // https://docs.zephyrproject.org/latest/build/sysbuild/index.html
		"--build-dir", workDir+"/build",
		workDir,
		"--",
		"-DNCS_TOOLCHAIN_VERSION=NONE",
		// FIXME: concat path in a os-appropriate way.
		fmt.Sprintf("-DCONF_FILE=%s/prj.conf", workDir),
		fmt.Sprintf("-DDTC_OVERLAY_FILE=%s/app.overlay", workDir),
	)

	toolchainsPath := device.General.GetToochainsPath()
	if err := build.Run(ctx, runner.WithToolchainPath(toolchainsPath.ToolchainPath, toolchainsPath.SDKPath)); err != nil {
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

func getWorkdir(cmd *cli.Command) (string, error) {
	workDir, err := filepath.Abs(cmd.String("workdir"))
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

func getConfigFile(cmd *cli.Command) string {
	if cmd.IsSet(filenameArg) {
		return cmd.String(filenameArg)
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
