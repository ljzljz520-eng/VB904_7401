package store

import (
	"sort"
	"sync"

	"example.com/videolab/internal/video"
)

type MemoryGradeRepository struct {
	mu     sync.RWMutex
	grades map[video.Scene]video.Grade
}

func NewMemoryGradeRepository() *MemoryGradeRepository {
	return &MemoryGradeRepository{grades: make(map[video.Scene]video.Grade)}
}

func (repository *MemoryGradeRepository) Save(grade video.Grade) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	repository.grades[grade.Scene] = grade
}

func (repository *MemoryGradeRepository) List() []video.Grade {
	repository.mu.RLock()
	defer repository.mu.RUnlock()
	grades := make([]video.Grade, 0, len(repository.grades))
	for _, grade := range repository.grades {
		grades = append(grades, grade)
	}
	sort.Slice(grades, func(left, right int) bool {
		return grades[left].Scene < grades[right].Scene
	})
	return grades
}
