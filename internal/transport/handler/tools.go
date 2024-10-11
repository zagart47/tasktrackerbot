package handler

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"tasktrackerbot/internal/entity"
	"time"
)

type Message struct {
	TaskId    string
	Text      string
	CreatedAt time.Time
	Reminder  time.Duration
	Complete  bool
}

func (m Message) String() string {
	status := "Нет"
	if m.Complete {
		status = "Да"
	}
	return fmt.Sprintf(
		"Номер задачи: %s\nОписание: %s\nДата добавления: %s\nДата исполнения: %s\nВыполнено: %s",
		m.TaskId, m.Text, m.CreatedAt.Format("15:04:05 02.01.2006"), m.Reminder, status)
}

func (b *Bot) sendReminder(task entity.Task) error {
	_, err := b.Send(tele.ChatID(task.UserID), fmt.Sprintf("Напоминание: %s", task.Text))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) getMainMenuKeyboard() *tele.ReplyMarkup {
	var (
		keyboard = &tele.ReplyMarkup{}
		btnList  = keyboard.Text("Мои задачи")
	)
	keyboard.ResizeKeyboard = true
	keyboard.Reply(
		keyboard.Row(btnList),
	)
	return keyboard
}
