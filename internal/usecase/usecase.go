package usecase

import (
	"context"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/repository"
)

type Tasks interface {
	AddTask(ctx context.Context, userID int64, taskText string, reminder entity.ReminderDuration) (entity.Task, error)
	GetTask(ctx context.Context, ID int64) (entity.Task, error)
	MarkReminderSent(ctx context.Context, id int64) error
	GetPendingReminders(ctx context.Context) ([]entity.Task, error)
}

type Usecases struct {
	Tasks Tasks
}

func NewUsecases(repo repository.Repositories) Usecases {
	taskUsecase := NewTaskUsecase(repo)
	return Usecases{
		Tasks: taskUsecase,
	}
}
