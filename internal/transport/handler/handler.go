package handler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"tasktrackerbot/pkg/remind"
	"time"

	tele "gopkg.in/telebot.v3"
	"tasktrackerbot/config"
	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/transport"
)

type Bot struct {
	*transport.BotService
}

func NewHandler(bot transport.BotService) Bot {
	return Bot{
		&bot,
	}
}

func (b *Bot) Start(c tele.Context) error {
	return c.Reply("hello")
}

func (b *Bot) HandleText(c tele.Context) error {
	// Сохраняем последнее сообщение от пользователя
	userID := c.Sender().ID
	lastMessages[userID] = c.Text()
	return c.Reply("Hello")
}

func checkCommand(command string, botName string) []string {
	// Проверяем, не соответствует ли команда формату @имя_бота ctrl 5d
	re := regexp.MustCompile(fmt.Sprintf(`@%s ctrl (\d+)([hdwms])`, botName))
	matches := re.FindStringSubmatch(command)
	if len(matches) == 3 {
		return matches
	}
	return nil
}

var lastMessages = map[int64]string{}

func (b *Bot) HandleCtrlCommand(c tele.Context) error {
	command := c.Text()
	botName := c.Bot().Me.Username

	matches := checkCommand(command, botName)
	id := c.Sender().ID
	if matches == nil {
		lastMessages[id] = command
		return nil
	}
	// Получаем последнее сообщение от пользователя
	text, ok := lastMessages[id]
	if !ok {
		return c.Reply("Не удалось найти задачу. Убедитесь, что вы отправили сообщение перед командой.")
	}

	// Извлекаем интервал и продолжительность
	interval, _ := strconv.Atoi(matches[1])
	duration := matches[2]

	// Создаем структуру ReminderDuration
	reminder := entity.ReminderDuration{
		Value: interval,
		Unit:  duration,
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	// Добавляем задачу с использованием сервиса
	task, err := b.Services.Tasks.AddTask(ctx, id, text, reminder)
	if err != nil {
		return c.Reply(fmt.Sprintf("Ошибка при добавлении задачи: %s", err.Error()))
	}

	// Планируем напоминание
	reminderTime := remind.CalculateReminderTime(reminder.Unit, reminder.Value)
	durationUntilReminder := reminderTime.Sub(time.Now())
	time.AfterFunc(durationUntilReminder, func() {
		b.sendReminder(task)
	})

	// Отправляем ответ пользователю
	response := fmt.Sprintf("#Задача# %s принята. Напомню о ней через %d%s.", text, interval, duration)
	return c.Reply(response)
}

func (b *Bot) sendReminder(task entity.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	err := b.Services.Tasks.MarkReminderSent(ctx, task.ID)
	if err != nil {
		fmt.Printf("Ошибка при отметке напоминания как отправленного: %s\n", err.Error())
		return
	}
	b.Send(tele.ChatID(task.UserID), fmt.Sprintf("Напоминание: %s", task.Text))
}

func (b *Bot) processReminders() {
	for {
		time.Sleep(time.Minute)
		ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
		tasks, err := b.Services.Tasks.GetPendingReminders(ctx)
		if err != nil {
			fmt.Printf("Ошибка при получении ожидающих напоминаний: %s\n", err.Error())
			continue
		}

		for _, task := range tasks {
			durationUntilReminder := task.ReminderTime.Sub(time.Now())
			if durationUntilReminder <= 0 {
				b.sendReminder(task)
			} else {
				time.AfterFunc(durationUntilReminder, func() {
					b.sendReminder(task)
				})
			}
		}
		cancel()
	}
}

func (b *Bot) restoreReminders() {
	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	tasks, err := b.Services.Tasks.GetPendingReminders(ctx)
	if err != nil {
		fmt.Printf("Ошибка при восстановлении напоминаний: %s\n", err.Error())
		return
	}

	for _, task := range tasks {
		durationUntilReminder := task.ReminderTime.Sub(time.Now())
		time.AfterFunc(durationUntilReminder, func() {
			b.sendReminder(task)
		})
	}
}

func (b *Bot) InitHandlers() {
	b.Bot.Handle(tele.OnText, b.HandleText)
	b.Bot.Handle("/start", b.Start)
	b.Bot.Handle(tele.OnText, b.HandleCtrlCommand) // Добавляем обработчик команды
	go b.processReminders()
	b.restoreReminders()
}
