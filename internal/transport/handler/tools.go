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
	status := "Нет"
	if m.Complete {
		status = "Да"
	}
	return fmt.Sprintf(
		"<b>Номер задачи:</b> %s\n<b>Описание:</b> %s\n<b>Дата добавления:</b> %s\n<b>Дата исполнения:</b> %s\n<b>Выполнено:</b> %s",
		m.TaskId, m.Text, m.CreatedAt.Format("15:04:05 02.01.2006"), m.Expiration.Format("15:04:05 02.01.2006"), status)
}

func (b *Bot) sendReminder(task entity.Task) error {
	_, err := b.Send(tele.ChatID(task.UserID), fmt.Sprintf("Напоминание: %s", task.Text))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) getMainMenuKeyboard() *tele.ReplyMarkup {
	keyboard := &tele.ReplyMarkup{}
	btnList := keyboard.Text("Мои задачи")
	keyboard.ResizeKeyboard = true
	keyboard.Reply(
		keyboard.Row(btnList),
	)
	return keyboard
}

func (b *Bot) getKeyboardForStartMessaging() *tele.ReplyMarkup {
	keyboard := &tele.ReplyMarkup{ResizeKeyboard: true}
	btn := keyboard.Data("Старт", "start_messaging", "start")
	keyboard.Inline(keyboard.Row(btn))
	return keyboard
}
