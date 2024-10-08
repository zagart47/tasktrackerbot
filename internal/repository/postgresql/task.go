package postgresql

import (
	"context"
	"tasktrackerbot/internal/entity"
)

type TaskRepo struct {
	db Client
}

func NewTaskRepo(db Client) TaskRepo {
	return TaskRepo{db: db}
}

func (r *TaskRepo) SaveTask(ctx context.Context, task entity.Task) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO tasks (user_id, text, created_at, reminder_time, reminder_sent)
        VALUES ($1, $2, $3, $4, $5)
    `, task.UserID, task.Text, task.CreatedAt, task.ReminderTime, task.ReminderSent)
	return err
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	var task entity.Task
	err := r.db.QueryRow(ctx, `
        SELECT id, user_id, text, created_at, reminder_time, reminder_sent FROM tasks WHERE id = $1
    `, id).Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.ReminderTime, &task.ReminderSent)
	if err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (r *TaskRepo) GetTasksByUserID(ctx context.Context, userID int64) ([]entity.Task, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, text, created_at, reminder_time, reminder_sent FROM tasks WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.ReminderTime, &task.ReminderSent); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) GetPendingReminders(ctx context.Context) ([]entity.Task, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, text, created_at, reminder_time, reminder_sent FROM tasks WHERE reminder_time > NOW() AND reminder_sent = FALSE
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Text, &task.CreatedAt, &task.ReminderTime, &task.ReminderSent); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) MarkReminderSent(ctx context.Context, taskID int64) error {
	_, err := r.db.Exec(ctx, `
        UPDATE tasks SET reminder_sent = TRUE WHERE id = $1
    `, taskID)
	return err
}
