package storage

import (
	"context"

	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/storage/postgresql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Tasks interface {
	CreateTask(ctx context.Context, task entity.Task) (int64, error)
	GetTaskByID(ctx context.Context, id int64) (entity.Task, error)
	GetTasksByUserID(ctx context.Context, userID int64) ([]entity.Task, error)
	GetUnsentTasks(ctx context.Context) ([]entity.Task, error)
	MarkTaskAsSent(ctx context.Context, taskID int64) error
}

type Storage struct {
	Tasks Tasks
}

func NewStorages(db *pgxpool.Pool) Storage {
	repos := postgresql.NewTaskStorage(db)
	return Storage{Tasks: &repos}
}
