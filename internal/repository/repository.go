package repository

import (
	"context"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/repository/postgresql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tasks interface {
	SaveTask(ctx context.Context, task entity.Task) error
	GetTaskByID(ctx context.Context, id int64) (entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID int64) ([]entity.Task, error)
	GetPendingReminders(ctx context.Context) ([]entity.Task, error)
	MarkReminderSent(ctx context.Context, taskID int64) error
}

type Repositories struct {
	Tasks Tasks
}

func NewRepositories(db *pgxpool.Pool) Repositories {
	repos := postgresql.NewTaskRepo(db)
	return Repositories{Tasks: &repos}
}
