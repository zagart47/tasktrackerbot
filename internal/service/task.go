package service

import (
	"context"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/usecase"
)

type Tasks interface {
	AddTask(context context.Context, id int64, text string, reminder entity.ReminderDuration) (entity.Task, error)
	GetTask(context context.Context, id int64) (entity.Task, error)
	MarkReminderSent(context context.Context, id int64) error
	GetPendingReminders(context context.Context) ([]entity.Task, error)
}

type TaskService struct {
	usecase usecase.Usecases
}

func NewTaskService(usecases usecase.Usecases) *TaskService {
	return &TaskService{
		usecase: usecases,
	}
}

func (s TaskService) AddTask(context context.Context, id int64, text string, reminder entity.ReminderDuration) (entity.Task, error) {
	return s.usecase.Tasks.AddTask(context, id, text, reminder)
}

func (s TaskService) GetTask(context context.Context, id int64) (entity.Task, error) {
	return s.usecase.Tasks.GetTask(context, id)
}

func (s TaskService) MarkReminderSent(ctx context.Context, id int64) error {
	return s.usecase.Tasks.MarkReminderSent(ctx, id)
}

func (s TaskService) GetPendingReminders(ctx context.Context) ([]entity.Task, error) {
	return s.usecase.Tasks.GetPendingReminders(ctx)
}
