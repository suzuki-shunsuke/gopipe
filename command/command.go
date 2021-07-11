package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/suzuki-shunsuke/go-timeout/timeout"
	"github.com/suzuki-shunsuke/gopipe/gopipe"
)

type CommandParam struct {
	Timeout time.Duration
}

func Command(cmd *exec.Cmd, param *CommandParam, opts ...CommandOpt) gopipe.Action {
	return func(ctx context.Context, args *gopipe.Args) error {
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}
		if cmd.Stderr == nil {
			cmd.Stderr = os.Stderr
		}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			if err := opt(cmd); err != nil {
				return err
			}
		}
		runner := timeout.NewRunner(0)
		if param != nil {
			if param.Timeout > 0 {
				c, cancel := context.WithTimeout(ctx, param.Timeout)
				defer cancel()
				ctx = c
			}
		}
		if err := runner.Run(ctx, cmd); err != nil {
			return err
		}
		if code := cmd.ProcessState.ExitCode(); code != 0 {
			return fmt.Errorf("exit code = %d != 0", code)
		}
		return nil
	}
}

type CommandOpt func(cmd *exec.Cmd) error

func Env(name, value string) CommandOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Env = append(cmd.Env, name+"="+value)
		return nil
	}
}

func Envs(envs []string) CommandOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Env = append(cmd.Env, envs...)
		return nil
	}
}

func EnvMap(m map[string]string) CommandOpt {
	return func(cmd *exec.Cmd) error {
		for name, value := range m {
			cmd.Env = append(cmd.Env, name+"="+value)
		}
		return nil
	}
}

func Dir(dir string) CommandOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Dir = dir
		return nil
	}
}

func Stdout(w io.Writer) CommandOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Stdout = w
		return nil
	}
}

func Stderr(w io.Writer) CommandOpt {
	return func(cmd *exec.Cmd) error {
		cmd.Stderr = w
		return nil
	}
}
