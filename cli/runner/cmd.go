package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Cmd struct {
	command string
	args    []string
}

type CmdOpt func(c *exec.Cmd)

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
		panic("please call AddArg again instead to add multiple arguments (for now)")
	}

	return c
}

func (c *Cmd) Run(ctx context.Context, opts ...CmdOpt) error {
	cmd := exec.CommandContext(ctx, c.command, c.args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(cmd.Env, os.Environ()...)

	for _, opt := range opts {
		opt(cmd)
	}

	updateCurrentPath(cmd.Env)

	// Always look for executable, if we have custom PATH present.
	commandPath, err := exec.LookPath(c.command)
	if err != nil {
		return fmt.Errorf("lookup command %q path: %w", c.command, err)
	}

	cmd.Path = commandPath
	cmd.Err = nil // Reset, as we tried to update the path of the command

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command %q: %w", c.command, err)
	}

	return nil
}

func WithWorkDir(workDir string) CmdOpt {
	return func(c *exec.Cmd) {
		c.Dir = workDir
	}
}

// WithToolchainPath updates environment of command
// to inlcude necessary variables for building firmware.
func WithToolchainPath(ncsToolchainBase, zephyrBase string) CmdOpt {
	if ncsToolchainBase == "" || zephyrBase == "" {
		return func(c *exec.Cmd) {}
	}

	return WithEnvironment(extendEnv(ncsToolchainBase, zephyrBase)...)
}

func WithEnvironment(envVals ...string) CmdOpt {
	return func(c *exec.Cmd) {
		// Prepend env to try and take higher priority.
		c.Env = append(c.Env, envVals...)
	}
}

func extendEnv(ncsToolchainPath string, zephyrPath string) []string {
	// All of this paths are based on NRF Connect SDK v2.5.1 and might change in the future.

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

	ldLibraryPath := generateEnvArray(ncsToolchainPath, []string{
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

	return strings.Join(vals, string(filepath.ListSeparator))
}

func updateCurrentPath(envs []string) {
	var envPath string
	for _, env := range envs {
		if !strings.HasPrefix(env, "PATH=") {
			continue
		}

		parts := strings.SplitN(env, "=", 2)

		if envPath == "" {
			envPath = parts[1]
		} else {
			envPath += string(filepath.ListSeparator) + parts[1]
		}
	}

	if envPath != "" {
		os.Setenv("PATH", envPath)
	}
}
