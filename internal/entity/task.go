package entity

import "time"

type Task struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"user_id"`
	Text         string        `json:"text"`
	CreatedAt    time.Time     `json:"created_at"`
	Expiration   time.Time     `json:"expiration"`
	Duration     time.Duration `json:"duration"`
	ReminderSent bool          `json:"reminder_sent"`
	MsgID        int           `json:"msg_id"`
	ChatID       int64         `json:"chat_id"`
	Reminder     Reminder
}

type Reminder struct {
	Value int
	Unit  string
}
