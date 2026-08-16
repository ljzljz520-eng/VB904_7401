package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"example.com/videolab/internal/app"
	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/fixture"
	"example.com/videolab/internal/store"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: videolab <plan|execute|simulate-error>")
		return 2
	}

	project, filters, scores, err := fixture.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	var runner ffmpeg.Runner
	execute := true
	switch args[0] {
	case "plan":
		execute = false
		runner = successfulRunner{}
	case "execute":
		runner = ffmpeg.OSRunner{}
	case "simulate-error":
		runner = decodeFailureRunner{}
	default:
		fmt.Fprintln(os.Stderr, "usage: videolab <plan|execute|simulate-error>")
		return 2
	}

	repository := store.NewMemoryGradeRepository()
	service := app.NewService(ffmpeg.NewExecutor(runner), repository, scores, filters)
	page := service.Produce(context.Background(), project, execute)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(page); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if page.Error != "" {
		return 1
	}
	return 0
}

type successfulRunner struct{}

func (successfulRunner) Run(context.Context, string, ...string) ([]byte, error) {
	return nil, nil
}

type decodeFailureRunner struct{}

func (decodeFailureRunner) Run(_ context.Context, _ string, args ...string) ([]byte, error) {
	for _, arg := range args {
		if arg == "libx264" {
			return []byte("Invalid data found while decoding stream #0:0"), errors.New("exit status 183")
		}
	}
	return nil, nil
}
