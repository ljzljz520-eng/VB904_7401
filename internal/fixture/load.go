package fixture

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"example.com/videolab/internal/video"
)

//go:embed fixtures/*.json
var files embed.FS

type projectDocument struct {
	Name    string         `json:"name"`
	Clips   []clipDocument `json:"clips"`
	Filters []video.Filter `json:"filters"`
}

type clipDocument struct {
	ID        string      `json:"id"`
	Scene     video.Scene `json:"scene"`
	Path      string      `json:"path"`
	FrameTime string      `json:"frame_time"`
}

type scoreDocument struct {
	Scene  video.Scene `json:"scene"`
	Filter string      `json:"filter"`
	Score  int         `json:"score"`
}

func Load() (video.Project, []video.Filter, video.ScoreBook, error) {
	var document projectDocument
	if err := readJSON("fixtures/project.json", &document); err != nil {
		return video.Project{}, nil, nil, err
	}

	project := video.Project{Name: document.Name, Clips: make([]video.Clip, 0, len(document.Clips))}
	for _, entry := range document.Clips {
		frameTime, err := time.ParseDuration(entry.FrameTime)
		if err != nil {
			return video.Project{}, nil, nil, fmt.Errorf("parse frame time for %s: %w", entry.ID, err)
		}
		project.Clips = append(project.Clips, video.Clip{
			ID: entry.ID, Scene: entry.Scene, Path: entry.Path, FrameTime: frameTime,
		})
	}

	var scoreEntries []scoreDocument
	if err := readJSON("fixtures/scores.json", &scoreEntries); err != nil {
		return video.Project{}, nil, nil, err
	}
	scores := make(video.ScoreBook, len(scoreEntries))
	for _, entry := range scoreEntries {
		scores[video.ScoreKey{Scene: entry.Scene, FilterName: entry.Filter}] = entry.Score
	}
	return project, append([]video.Filter(nil), document.Filters...), scores, nil
}

func readJSON(path string, destination any) error {
	data, err := files.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read fixture %s: %w", path, err)
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return fmt.Errorf("decode fixture %s: %w", path, err)
	}
	return nil
}
