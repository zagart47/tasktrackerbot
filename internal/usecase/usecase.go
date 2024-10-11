package usecase

import (
	"context"
	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/storage"
	"tasktrackerbot/internal/storage/cache"
)

type Tasks interface {
	AddTask(ctx context.Context, task entity.Task) (int64, error)
	GetTaskByID(ctx context.Context, ID int64) (entity.Task, error)
	GetUnsentTasks(ctx context.Context) ([]entity.Task, error)
	GetTasksByUserID(ctx context.Context, userId int64) ([]entity.Task, error)
	MarkAsSent(ctx context.Context, id int64) error
	MakeTasksCache(ctx context.Context) error
}

type Usecases struct {
	Tasks Tasks
	Cache cache.Cacher
}

func NewUsecases(repo storage.Storage, mc cache.Cacher) Usecases {
	taskUsecase := NewTaskUsecase(repo, mc)
	return Usecases{
		Tasks: &taskUsecase,
	}
}
