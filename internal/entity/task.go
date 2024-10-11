package entity

import "time"

type Task struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	Text         string        `json:"text"`
	CreatedAt    time.Time     `json:"created_at"`
	Expiration   time.Time     `json:"expiration"`
	Duration     time.Duration `json:"duration"`
	Reminder     Reminder
	ReminderSent bool `json:"reminder_sent"`
}

type Reminder struct {
	Value int
	Unit  string
}
