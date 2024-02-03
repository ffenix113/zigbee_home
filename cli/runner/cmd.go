package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Cmd struct {
	command string
	args    []string
}

func NewCmd(command string, args ...string) *Cmd {
	return &Cmd{
		command: command,
		args:    args,
	}
}

func (c *Cmd) AddArg(arg string, value ...string) *Cmd {
	c.args = append(c.args, arg)

	switch len(value) {
	case 0:
	case 1:
		c.args = append(c.args, value[0])
	default:
		panic("please call AddArg again instead (for now)")
	}

	return c
}

func (c *Cmd) Run(ctx context.Context, workDir string) error {
	ncsToolchainPath := os.Getenv("NCS_TOOLCHAIN_BASE")
	zephyrPath := os.Getenv("ZEPHYR_BASE")

	var extendedEnv []string
	if ncsToolchainPath != "" && zephyrPath != "" {
		extendedEnv = c.extendEnv(ncsToolchainPath, zephyrPath)

		newPath := extendedEnv[0]
		pathParts := strings.SplitN(newPath, "=", 2)
		os.Setenv("PATH", pathParts[1])

		commandPath, err := exec.LookPath(c.command)
		if err != nil {
			return err
		}

		c.command = commandPath
	}

	cmd := exec.CommandContext(ctx, c.command, c.args...)

	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(extendedEnv, append(cmd.Env, os.Environ()...)...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command %q: %w", c.command, err)
	}

	return nil
}

func (c *Cmd) extendEnv(ncsToolchainPath string, zephyrPath string) []string {
	ncsCombinedPath := generateEnvArray(ncsToolchainPath, []string{
		"/usr/bin",
		"/usr/local/bin",
		"/opt/bin",
		"/opt/nanopb/generator-bin",
		"/opt/zephyr-sdk/aarch64-zephyr-elf/bin",
		"/opt/zephyr-sdk/x86_64-zephyr-elf/bin",
		"/opt/zephyr-sdk/arm-zephyr-eabi/bin",
	})

	envPath := os.Getenv("PATH")
	combinedPath := ncsCombinedPath
	if envPath != "" {
		combinedPath += ":" + envPath
	}

	pythonPath := generateEnvArray(ncsToolchainPath, []string{
		"/usr/local/lib/python3.8",
		"/usr/local/lib/python3.8/site-packages",
	})

	ldLibraryPath := generateEnvArray(ncsCombinedPath, []string{
		"/usr/lib",
		"/usr/lib/x86_64-linux-gnu",
		"/usr/local/lib",
	})

	return []string{
		"PATH=" + combinedPath,
		"ZEPHYR_BASE=" + zephyrPath,
		"ZEPHYR_SDK_INSTALL_DIR=" + path.Join(ncsToolchainPath, "/opt/zephyr-sdk"),
		"ZEPHYR_TOOLCHAIN_VARIANT=zephyr",
		"PYTHOHOME=" + path.Join(ncsCombinedPath, "/usr/local"),
		"PYTHONPATH=" + pythonPath,
		"LD_LIBRARY_PATH=" + ldLibraryPath,
	}
}

func generateEnvArray(prefix string, vals []string) string {
	for i := range vals {
		vals[i] = path.Join(prefix, vals[i])
	}

	return strings.Join(vals, ":")
}
