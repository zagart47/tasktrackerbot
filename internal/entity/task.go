package entity

import "time"

type Task struct {
	ID           int64
	UserID       int64
	Text         string
	Reminder     ReminderDuration
	CreatedAt    time.Time
	ReminderTime time.Time
	ReminderSent bool
}

type ReminderDuration struct {
	Value int
	Unit  string
}
