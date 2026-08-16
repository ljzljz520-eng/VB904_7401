package ffmpeg

import (
	"context"
	"errors"
	"os/exec"

	"example.com/videolab/internal/video"
)

type Runner interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

type OSRunner struct{}

func (OSRunner) Run(ctx context.Context, program string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, program, args...).CombinedOutput()
}

type Executor struct {
	runner Runner
}

func NewExecutor(runner Runner) *Executor {
	return &Executor{runner: runner}
}

func (executor *Executor) Execute(ctx context.Context, _ string, command video.Command) error {
	_, err := executor.runner.Run(ctx, command.Program, command.Args...)
	if err != nil {
		return errors.New("command failed")
	}
	return nil
}
