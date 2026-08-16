package app

import (
	"context"
	"fmt"
	"sort"

	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/video"
)

type Page struct {
	Title       string                  `json:"title"`
	Comparisons []video.FrameComparison `json:"comparisons,omitempty"`
	Schemes     []video.Grade           `json:"schemes,omitempty"`
	Renders     []video.RenderPlan      `json:"renders,omitempty"`
	Error       string                  `json:"error,omitempty"`
}

type Service struct {
	executor video.CommandExecutor
	grades   video.GradeRepository
	scores   video.ScoreBook
	filters  []video.Filter
}

func NewService(executor video.CommandExecutor, grades video.GradeRepository, scores video.ScoreBook, filters []video.Filter) *Service {
	return &Service{
		executor: executor,
		grades:   grades,
		scores:   scores,
		filters:  append([]video.Filter(nil), filters...),
	}
}

func (service *Service) Produce(ctx context.Context, project video.Project, execute bool) Page {
	page := Page{Title: project.Name}
	if err := service.validate(project); err != nil {
		page.Error = err.Error()
		return page
	}

	selected := make(map[video.Scene]video.Grade, len(project.Clips))
	for _, clip := range project.Clips {
		for _, filter := range service.filters {
			command, framePath := ffmpeg.FrameCommand(clip, filter)
			comparison := video.FrameComparison{
				Scene:       clip.Scene,
				VideoPath:   clip.Path,
				FilterName:  filter.Name,
				Score:       service.scores[video.ScoreKey{Scene: clip.Scene, FilterName: filter.Name}],
				FramePath:   framePath,
				Command:     command,
				WindowsLine: ffmpeg.WindowsLine(command),
			}
			page.Comparisons = append(page.Comparisons, comparison)
			if execute {
				if err := service.executor.Execute(ctx, clip.Path, command); err != nil {
					page.Error = err.Error()
					return page
				}
			}
			candidate := video.Grade{Scene: clip.Scene, FilterName: filter.Name, Graph: filter.Graph, Score: comparison.Score}
			current, present := selected[clip.Scene]
			if !present || better(candidate, current) {
				selected[clip.Scene] = candidate
			}
		}
	}

	for _, clip := range project.Clips {
		grade := selected[clip.Scene]
		service.grades.Save(grade)
		command, outputPath := ffmpeg.RenderCommand(clip, grade)
		plan := video.RenderPlan{
			Scene:       clip.Scene,
			VideoPath:   clip.Path,
			FilterName:  grade.FilterName,
			OutputPath:  outputPath,
			Command:     command,
			WindowsLine: ffmpeg.WindowsLine(command),
		}
		page.Renders = append(page.Renders, plan)
	}
	page.Schemes = service.grades.List()

	if execute {
		for _, plan := range page.Renders {
			if err := service.executor.Execute(ctx, plan.VideoPath, plan.Command); err != nil {
				page.Error = err.Error()
				return page
			}
		}
	}
	return page
}

func (service *Service) validate(project video.Project) error {
	if project.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if len(project.Clips) == 0 {
		return fmt.Errorf("at least one clip is required")
	}
	if len(service.filters) == 0 {
		return fmt.Errorf("at least one filter is required")
	}
	for _, clip := range project.Clips {
		if clip.ID == "" || clip.Path == "" {
			return fmt.Errorf("clip id and path are required")
		}
		for _, filter := range service.filters {
			key := video.ScoreKey{Scene: clip.Scene, FilterName: filter.Name}
			if _, present := service.scores[key]; !present {
				return fmt.Errorf("missing score for %s/%s", clip.Scene, filter.Name)
			}
		}
	}
	return nil
}

func better(candidate, current video.Grade) bool {
	if candidate.Score != current.Score {
		return candidate.Score > current.Score
	}
	names := []string{candidate.FilterName, current.FilterName}
	sort.Strings(names)
	return candidate.FilterName == names[0]
}
