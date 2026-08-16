package ffmpeg

import (
	"strconv"
	"strings"
	"time"

	"example.com/videolab/internal/video"
)

func FrameCommand(clip video.Clip, filter video.Filter) (video.Command, string) {
	output := "frames\\" + clip.ID + "-" + filter.Name + ".jpg"
	command := video.Command{
		Program: "ffmpeg",
		Args: []string{
			"-hide_banner",
			"-y",
			"-ss",
			formatTimestamp(clip.FrameTime),
			"-i",
			clip.Path,
			"-vf",
			filter.Graph,
			"-frames:v",
			"1",
			output,
		},
	}
	return command, output
}

func RenderCommand(clip video.Clip, grade video.Grade) (video.Command, string) {
	output := "renders\\" + clip.ID + "-graded.mp4"
	command := video.Command{
		Program: "ffmpeg",
		Args: []string{
			"-hide_banner",
			"-y",
			"-i",
			clip.Path,
			"-vf",
			grade.Graph,
			"-c:v",
			"libx264",
			"-c:a",
			"copy",
			output,
		},
	}
	return command, output
}

func WindowsLine(command video.Command) string {
	parts := make([]string, 0, len(command.Args)+1)
	parts = append(parts, quoteWindowsArg(command.Program))
	for _, arg := range command.Args {
		parts = append(parts, quoteWindowsArg(arg))
	}
	return strings.Join(parts, " ")
}

func formatTimestamp(value time.Duration) string {
	hours := int(value / time.Hour)
	value -= time.Duration(hours) * time.Hour
	minutes := int(value / time.Minute)
	value -= time.Duration(minutes) * time.Minute
	seconds := float64(value) / float64(time.Second)
	return twoDigits(hours) + ":" + twoDigits(minutes) + ":" + formatSeconds(seconds)
}

func twoDigits(value int) string {
	if value < 10 {
		return "0" + strconv.Itoa(value)
	}
	return strconv.Itoa(value)
}

func formatSeconds(value float64) string {
	formatted := strconv.FormatFloat(value, 'f', 3, 64)
	if value < 10 {
		return "0" + formatted
	}
	return formatted
}

func quoteWindowsArg(value string) string {
	if value != "" && !strings.ContainsAny(value, " \t\n\v\"") {
		return value
	}
	var builder strings.Builder
	builder.WriteByte('"')
	backslashes := 0
	for _, character := range value {
		switch character {
		case '\\':
			backslashes++
		case '"':
			builder.WriteString(strings.Repeat("\\", backslashes*2+1))
			builder.WriteRune(character)
			backslashes = 0
		default:
			builder.WriteString(strings.Repeat("\\", backslashes))
			builder.WriteRune(character)
			backslashes = 0
		}
	}
	builder.WriteString(strings.Repeat("\\", backslashes*2))
	builder.WriteByte('"')
	return builder.String()
}
