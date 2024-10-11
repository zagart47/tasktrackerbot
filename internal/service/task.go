package service

import (
	"context"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/usecase"
)

type Tasks interface {
	AddTask(context context.Context, task entity.Task) (int64, error)
	GetTaskByID(context context.Context, id int64) (entity.Task, error)
	GetTasksByUserID(context context.Context, userId int64) ([]entity.Task, error)
	GetUnsentTasks(context context.Context) ([]entity.Task, error)
	MakeTasksCache(context context.Context) error
	MarkAsSent(context context.Context, id int64) error
}

type TaskService struct {
	usecase usecase.Usecases
}

func (s TaskService) MakeTasksCache(context context.Context) error {
	return s.usecase.Tasks.MakeTasksCache(context)
}

func NewTaskService(usecase usecase.Usecases) *TaskService {
	return &TaskService{usecase: usecase}
}

func (s TaskService) AddTask(context context.Context, task entity.Task) (int64, error) {
	return s.usecase.Tasks.AddTask(context, task)
}

func (s TaskService) GetTaskByID(context context.Context, id int64) (entity.Task, error) {
	return s.usecase.Tasks.GetTaskByID(context, id)
}

func (s TaskService) GetUnsentTasks(context context.Context) ([]entity.Task, error) {
	return s.usecase.Tasks.GetUnsentTasks(context)
}

func (s TaskService) MarkAsSent(ctx context.Context, id int64) error {
	return s.usecase.Tasks.MarkAsSent(ctx, id)
}

func (s TaskService) GetTasksByUserID(context context.Context, userId int64) ([]entity.Task, error) {
	return s.usecase.Tasks.GetTasksByUserID(context, userId)
}
