package ffmpeg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

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

func (executor *Executor) Execute(ctx context.Context, videoPath string, command video.Command) error {
	output, err := executor.runner.Run(ctx, command.Program, command.Args...)
	if err != nil {
		detail := strings.TrimSpace(string(output))
		if detail == "" {
			return fmt.Errorf("ffmpeg failed for \"%s\": %w", videoPath, err)
		}
		return fmt.Errorf("ffmpeg failed for \"%s\": %w: %s", videoPath, err, detail)
	}
	return nil
}
