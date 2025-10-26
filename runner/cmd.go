package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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
func WithToolchainPath(ncsToolchainBase, sdkBase string) CmdOpt {
	// For now check that we don't want to setup env here,
	// and move it to CLI ASAP.
	// This could be useful if run inside environment that
	// is already set up properly.
	if noSetupEnv() || ncsToolchainBase == "" || sdkBase == "" {
		log.Println("environment will not be prepared because either one of the paths is empty, or requested not to")

		return func(c *exec.Cmd) {}
	}

	return WithEnvironment(extendEnv(ncsToolchainBase, sdkBase)...)
}

func WithEnvironment(envVals ...string) CmdOpt {
	return func(c *exec.Cmd) {
		// Prepend env to try and take higher priority.
		c.Env = append(c.Env, envVals...)
	}
}

func extendEnv(ncsToolchainPath string, sdkPath string) []string {
	envFilePath := filepath.Join(ncsToolchainPath, "environment.json")
	envFile, err := os.Open(envFilePath)
	if err != nil {
		log.Printf("error opening environment.json file at %q: %v", envFilePath, err)
		return nil
	}
	defer envFile.Close()

	var envConfig struct {
		EnvVars []struct {
			Type                   string   `json:"type"`
			Key                    string   `json:"key"`
			Values                 []string `json:"values"`
			Value                  string   `json:"value"`
			ExistingValueTreatment string   `json:"existing_value_treatment"`
		} `json:"env_vars"`
	}

	if err := json.NewDecoder(envFile).Decode(&envConfig); err != nil {
		log.Printf("error decoding environment.json file: %v", err)
		return nil
	}

	envVars := make(map[string]string)
	for _, envVar := range envConfig.EnvVars {
		switch envVar.Type {
		case "relative_paths":
			paths := generateEnvArray(ncsToolchainPath, envVar.Values)
			if envVar.ExistingValueTreatment == "prepend_to" {
				existingValue := os.Getenv(envVar.Key)
				if existingValue != "" {
					paths += string(os.PathListSeparator) + existingValue
				}
			}
			envVars[envVar.Key] = paths
		case "string":
			envVars[envVar.Key] = envVar.Value
		}
	}

	ncsCombinedPath := envVars["PATH"]

	envPath := os.Getenv("PATH")
	combinedPath := ncsCombinedPath

	if envPath != "" {
		combinedPath += string(os.PathListSeparator) + envPath
	}

	// This is Linux specific, so may not work correctly on Windows.
	ldLibraryPath := generateEnvArray(ncsToolchainPath, []string{
		"/usr/lib",
		"/usr/lib/x86_64-linux-gnu",
		"/usr/local/lib",
	})

	return []string{
		"PATH=" + combinedPath,
		"ZEPHYR_BASE=" + sdkPath,
		"ZEPHYR_SDK_INSTALL_DIR=" + filepath.Join(ncsToolchainPath, "opt", "zephyr-sdk"),
		"ZEPHYR_TOOLCHAIN_VARIANT=zephyr",
		"LD_LIBRARY_PATH=" + ldLibraryPath,
	}
}

func generateEnvArray(prefix string, vals []string) string {
	for i := range vals {
		vals[i] = filepath.Join(prefix, vals[i])
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

func noSetupEnv() bool {
	_, ok := os.LookupEnv("NO_SETUP_ENV")

	return ok
}
