package ffmpeg

import (
	"strings"
	"testing"
	"time"

	"example.com/videolab/internal/video"
)

func TestFrameCommandQuotesWindowsPaths(t *testing.T) {
	clip := video.Clip{
		ID: "sea", Scene: video.SceneSeaside,
		Path: "C:\\Video Projects\\Beach Takes\\sea breeze.mp4", FrameTime: 5250 * time.Millisecond,
	}
	filter := video.Filter{Name: "vivid", Graph: "eq=contrast=1.08:saturation=1.22"}
	command, output := FrameCommand(clip, filter)
	line := WindowsLine(command)
	if output != "frames\\sea-vivid.jpg" {
		t.Fatalf("output = %q", output)
	}
	if !strings.Contains(line, `-i "C:\Video Projects\Beach Takes\sea breeze.mp4"`) {
		t.Fatalf("command does not quote input: %s", line)
	}
	if !strings.Contains(line, "-ss 00:00:05.250") {
		t.Fatalf("command timestamp = %s", line)
	}
}
