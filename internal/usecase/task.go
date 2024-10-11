package usecase

import (
	"context"
	"strconv"
	"time"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/storage"
	"tasktrackerbot/internal/storage/cache"
)

type TaskUsecase struct {
	repo storage.Storage
	mc   cache.Cacher
}

func NewTaskUsecase(repo storage.Storage, mc cache.Cacher) TaskUsecase {
	return TaskUsecase{
		repo: repo,
		mc:   mc,
	}
}

func (u *TaskUsecase) AddTask(ctx context.Context, task entity.Task) (int64, error) {
	var err error
	task.ID, err = u.repo.Tasks.CreateTask(ctx, task)
	if err != nil {
		return 0, err
	}
	err = u.mc.Set(task)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (u *TaskUsecase) GetTaskByID(ctx context.Context, ID int64) (entity.Task, error) {
	return u.repo.Tasks.GetTaskByID(ctx, ID)
}

func (u *TaskUsecase) GetUnsentTasks(ctx context.Context) ([]entity.Task, error) {
	tasks, err := u.mc.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var expiredTasks []entity.Task
	for _, task := range tasks {
		if task.Expiration.Before(time.Now()) {
			expiredTasks = append(expiredTasks, task)
		}
	}
	return expiredTasks, nil
}

func (u *TaskUsecase) MarkAsSent(ctx context.Context, taskId int64) error {
	err := u.mc.Delete(strconv.FormatInt(taskId, 10))
	if err != nil {
		return err
	}
	return u.repo.Tasks.MarkTaskAsSent(ctx, taskId)
}

func (u *TaskUsecase) MakeTasksCache(ctx context.Context) error {
	tasks, err := u.repo.Tasks.GetUnsentTasks(ctx)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		err = u.mc.Set(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *TaskUsecase) GetTasksByUserID(ctx context.Context, userId int64) ([]entity.Task, error) {
	return u.repo.Tasks.GetTasksByUserID(ctx, userId)
}
