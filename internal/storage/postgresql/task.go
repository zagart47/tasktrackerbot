package postgresql

import (
	"context"
	"tasktrackerbot/internal/entity"
)

type TaskStorage struct {
	db Client
}

func NewTaskStorage(db Client) TaskStorage {
	return TaskStorage{db: db}
}

func (s *TaskStorage) CreateTask(ctx context.Context, task entity.Task) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, `
        INSERT INTO tasks (user_id, text, created_at, expiration, duration)
        VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
    `, task.UserID, task.Text, task.CreatedAt, task.Expiration, task.Duration).Scan(&id)
	return id, err
}

func (s *TaskStorage) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	var task entity.Task
	err := s.db.QueryRow(ctx, `
        SELECT id, user_id, text, created_at, expiration, duration, reminder_sent FROM tasks WHERE id = $1
    `, id).Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.Expiration, &task.Duration, &task.ReminderSent)
	if err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (s *TaskStorage) GetTasksByUserID(ctx context.Context, userID int64) ([]entity.Task, error) {
	rows, err := s.db.Query(ctx, `
        SELECT id, user_id, text, created_at, expiration, duration, reminder_sent FROM tasks WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.Expiration, &task.Duration, &task.ReminderSent); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStorage) GetUnsentTasks(ctx context.Context) ([]entity.Task, error) {
	rows, err := s.db.Query(ctx, `
        SELECT id, user_id, text, created_at, expiration, duration, reminder_sent FROM tasks WHERE reminder_sent = FALSE
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.Expiration, &task.Duration, &task.ReminderSent); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStorage) MarkTaskAsSent(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx, `
        UPDATE tasks SET reminder_sent = TRUE WHERE id = $1
    `, id)
	return err
}
