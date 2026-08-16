package video

import "time"

type Scene string

const (
	SceneSeaside Scene = "seaside"
	SceneCity    Scene = "city"
	SceneNight   Scene = "night"
)

type Clip struct {
	ID        string        `json:"id"`
	Scene     Scene         `json:"scene"`
	Path      string        `json:"path"`
	FrameTime time.Duration `json:"-"`
}

type Project struct {
	Name  string `json:"name"`
	Clips []Clip `json:"clips"`
}

type Filter struct {
	Name  string `json:"name"`
	Graph string `json:"graph"`
}

type Grade struct {
	Scene      Scene  `json:"scene"`
	FilterName string `json:"filter_name"`
	Graph      string `json:"graph"`
	Score      int    `json:"score"`
}

type ScoreKey struct {
	Scene      Scene
	FilterName string
}

type ScoreBook map[ScoreKey]int

type Command struct {
	Program string   `json:"program"`
	Args    []string `json:"args"`
}

type FrameComparison struct {
	Scene       Scene   `json:"scene"`
	VideoPath   string  `json:"video_path"`
	FilterName  string  `json:"filter_name"`
	Score       int     `json:"score"`
	FramePath   string  `json:"frame_path"`
	Command     Command `json:"-"`
	WindowsLine string  `json:"windows_command"`
}

type RenderPlan struct {
	Scene       Scene   `json:"scene"`
	VideoPath   string  `json:"video_path"`
	FilterName  string  `json:"filter_name"`
	OutputPath  string  `json:"output_path"`
	Command     Command `json:"-"`
	WindowsLine string  `json:"windows_command"`
}
