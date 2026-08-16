package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"example.com/videolab/internal/app"
	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/fixture"
	"example.com/videolab/internal/store"
)

func TestProductionPlansFramesGradesAndWindowsCommands(t *testing.T) {
	project, filters, scores, err := fixture.Load()
	if err != nil {
		t.Fatal(err)
	}
	repository := store.NewMemoryGradeRepository()
	runner := &recordingRunner{}
	service := app.NewService(ffmpeg.NewExecutor(runner), repository, scores, filters)
	page := service.Produce(context.Background(), project, true)
	if page.Error != "" {
		t.Fatalf("page error = %q", page.Error)
	}
	if len(page.Comparisons) != 9 {
		t.Fatalf("comparisons = %d", len(page.Comparisons))
	}
	if len(page.Schemes) != 3 || len(repository.List()) != 3 {
		t.Fatalf("schemes = %d stored = %d", len(page.Schemes), len(repository.List()))
	}
	if len(page.Renders) != 3 || len(runner.calls) != 12 {
		t.Fatalf("renders = %d calls = %d", len(page.Renders), len(runner.calls))
	}
	wantFilters := []string{"vivid", "cinematic", "vivid"}
	for index, render := range page.Renders {
		if render.FilterName != wantFilters[index] {
			t.Fatalf("render %d filter = %q", index, render.FilterName)
		}
		if !strings.Contains(render.WindowsLine, `-i "C:\Video Projects\`) {
			t.Fatalf("render %d command = %s", index, render.WindowsLine)
		}
	}
}

func TestProductionPageReportsDecodeFailure(t *testing.T) {
	project, filters, scores, err := fixture.Load()
	if err != nil {
		t.Fatal(err)
	}
	repository := store.NewMemoryGradeRepository()
	runner := &recordingRunner{failRender: true}
	service := app.NewService(ffmpeg.NewExecutor(runner), repository, scores, filters)
	page := service.Produce(context.Background(), project, true)
	if len(repository.List()) != 3 {
		t.Fatalf("stored schemes = %d", len(repository.List()))
	}
	if len(runner.calls) != 10 {
		t.Fatalf("calls before failure = %d", len(runner.calls))
	}
	if !strings.Contains(page.Error, project.Clips[0].Path) {
		t.Fatalf("page error omits video path: %q", page.Error)
	}
	if !strings.Contains(page.Error, "Invalid data found while decoding stream #0:0") {
		t.Fatalf("page error omits ffmpeg detail: %q", page.Error)
	}
}

type recordingRunner struct {
	calls      []recordedCall
	failRender bool
}

type recordedCall struct {
	program string
	args    []string
}

func (runner *recordingRunner) Run(_ context.Context, program string, args ...string) ([]byte, error) {
	runner.calls = append(runner.calls, recordedCall{program: program, args: append([]string(nil), args...)})
	if runner.failRender && contains(args, "libx264") {
		return []byte("Invalid data found while decoding stream #0:0"), errors.New("exit status 183")
	}
	return nil, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
