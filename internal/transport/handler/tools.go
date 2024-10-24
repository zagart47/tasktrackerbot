package handler

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"tasktrackerbot/internal/entity"
	"time"
)

type Message struct {
	TaskId     string
	Text       string
	CreatedAt  time.Time
	Expiration time.Time
	Complete   bool
}

func (m Message) String() string {
	return fmt.Sprintf(
		"<b>Номер задачи:</b> %s\n<b>Описание:</b> %s\n<b>Дата напоминания:</b> %s",
		m.TaskId, m.Text, m.Expiration.Format("15:04:05 02.01.2006"))
}

func (b *Bot) sendReminder(task entity.Task) error {
	_, err := b.Send(tele.ChatID(task.UserID), fmt.Sprintf("Напоминание: %s", task.Text))
	if err != nil {
		return err
	}
	return nil
}
