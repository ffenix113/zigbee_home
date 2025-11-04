package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Cmd struct {
	command string
	args    []string
}

// RunContext holds pre-defined context for all commands
// that will be executed by this context.
//
// Adding options to context will create a new context.
// Reciever will not be modified.
type RunContext struct {
	// While this means that for each command we
	// would need to re-construct execution environment
	// it is still good enough as initial solution.
	opts []CmdOpt
}

type CmdOpt func(c *exec.Cmd)

func NewRunContext(opts ...CmdOpt) RunContext {
	return RunContext{opts}
}

func (c RunContext) AddOpts(opts ...CmdOpt) RunContext {
	return RunContext{append(c.opts, opts...)}
}

func NewCmd(command string, args ...string) *Cmd {
	return &Cmd{
		command: command,
		args:    args,
	}
}

func (c *RunContext) Run(ctx context.Context, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(cmd.Env, os.Environ()...)

	for _, opt := range c.opts {
		opt(cmd)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command %q: %w", command, err)
	}

	return nil
}

func (c *Cmd) Run(ctx context.Context, opts ...CmdOpt) error {
	cmd := exec.CommandContext(ctx, c.command, c.args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(cmd.Env, os.Environ()...)

	for _, opt := range opts {
		opt(cmd)
	}

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
