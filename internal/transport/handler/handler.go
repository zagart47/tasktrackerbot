package handler

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"tasktrackerbot/pkg/remind"
	"time"

	"tasktrackerbot/config"
	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/transport"

	tele "gopkg.in/telebot.v3"
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
	return b.sendMessage(c, "Привет! Я бот для управления задачами.", b.getMainMenuKeyboard())
}

func (b *Bot) sendMessage(c tele.Context, text string, keyboard *tele.ReplyMarkup) error {
	return c.Send(text, &tele.SendOptions{ReplyMarkup: keyboard, ParseMode: tele.ModeHTML})
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
	msg, ok := lastMessages[id]
	if !ok {
		return c.Reply("Не удалось найти задачу. Убедитесь, что вы отправили сообщение перед командой.")
	}

	// Извлекаем интервал и продолжительность
	duration, _ := strconv.Atoi(matches[1])
	timeFormat := matches[2]

	err := remind.ValidateReminderDuration(timeFormat, duration)
	if err != nil {
		return c.Reply("Неверный формат срока напоминания")
	}

	taskDuration, timeUnit := remind.CalculateReminderTime(timeFormat, duration)

	task := entity.Task{
		UserID:       id,
		Text:         msg,
		CreatedAt:    time.Now(),
		Expiration:   time.Now().Add(taskDuration),
		Duration:     taskDuration,
		ReminderSent: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	// Добавляем задачу с использованием сервиса
	task.ID, err = b.Services.Tasks.AddTask(ctx, task)
	if err != nil {
		return c.Reply(fmt.Sprintf("Ошибка при добавлении задачи: %s", err.Error()))
	}

	// Отправляем ответ пользователю

	text := fmt.Sprintf("#Задача# \"%s\" принята. Присвоил ей № <b>%v</b>. Напомню о ней через <b>%d %s</b>.", task.Text, task.ID, duration, timeUnit)
	return c.Reply(text, &tele.SendOptions{ParseMode: tele.ModeHTML})
}

func (b *Bot) InitHandlers() {
	b.Bot.Handle("/start", b.Start)
	b.Bot.Handle(tele.OnText, b.HandleCtrlCommand) // Добавляем обработчик команды
	b.Bot.Handle(&tele.ReplyButton{Text: "Мои задачи"}, b.MyTasksHandler)
	go b.StartTasksSending()
}

func (b *Bot) StartTasksSending() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
		tasks, err := b.Services.Tasks.GetUnsentTasks(ctx)
		cancel()
		if err != nil {
			log.Println(err)
		}
		for _, task := range tasks {
			ctx, cancel = context.WithTimeout(context.Background(), config.Configs.Timeout)
			err = b.sendReminder(task)
			if err != nil {
				log.Println(err)
				continue
			}
			err = b.Services.Tasks.MarkAsSent(ctx, task.ID)
			if err != nil {
				log.Println(err)
				continue
			}
			cancel()
		}
		time.Sleep(time.Second * 1)
	}
}

func (b *Bot) MyTasksHandler(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	tasks, err := b.Services.Tasks.GetTasksByUserID(ctx, c.Sender().ID)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		m := Message{
			TaskId:    strconv.FormatInt(task.ID, 10),
			Text:      task.Text,
			CreatedAt: task.CreatedAt,
			Reminder:  time.Duration(task.Reminder.Value),
			Complete:  task.ReminderSent,
		}
		text := m.String()
		b.sendMessage(c, text, b.getMainMenuKeyboard())
	}
	return nil
}
