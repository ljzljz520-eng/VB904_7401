package video

import "context"

type CommandExecutor interface {
	Execute(context.Context, string, Command) error
}

type GradeRepository interface {
	Save(Grade)
	List() []Grade
}
