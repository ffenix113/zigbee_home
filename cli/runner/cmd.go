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
	cmd := exec.CommandContext(ctx, c.command, c.args...)

	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = append(cmd.Env, os.Environ()...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command %q: %w", c.command, err)
	}

	return nil
}
