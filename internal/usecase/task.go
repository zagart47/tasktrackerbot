package usecase

import (
	"context"
	"time"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/repository"
	"tasktrackerbot/pkg/remind"
)

type TaskUsecase struct {
	repo repository.Repositories
}

func NewTaskUsecase(repo repository.Repositories) TaskUsecase {
	return TaskUsecase{repo: repo}
}

func (u TaskUsecase) AddTask(ctx context.Context, userID int64, taskText string, duration entity.ReminderDuration) (entity.Task, error) {
	if err := remind.ValidateReminderDuration(duration.Unit, duration.Value); err != nil {
		return entity.Task{}, err
	}

	reminderTime := remind.CalculateReminderTime(duration.Unit, duration.Value)
	task := entity.Task{
		UserID:       userID,
		Text:         taskText,
		ReminderTime: reminderTime,
		CreatedAt:    time.Now(),
	}
	err := u.repo.Tasks.SaveTask(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (u TaskUsecase) GetTask(ctx context.Context, ID int64) (entity.Task, error) {
	return u.repo.Tasks.GetTaskByID(ctx, ID)
}

func (u TaskUsecase) MarkReminderSent(ctx context.Context, id int64) error {
	return u.repo.Tasks.MarkReminderSent(ctx, id)
}

func (u TaskUsecase) GetPendingReminders(ctx context.Context) ([]entity.Task, error) {
	return u.repo.Tasks.GetPendingReminders(ctx)
}
