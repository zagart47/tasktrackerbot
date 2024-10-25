package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"tasktrackerbot/config"
	"tasktrackerbot/internal/entity"
	"tasktrackerbot/internal/transport"
	"tasktrackerbot/pkg/remind"

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
	return b.sendMessage(c, "Если ты видишь это сообщение, значит всё идет правильно. Теперь я буду запоминать твои задачи и напоминать тебе о них здесь.")
}

func (b *Bot) sendMessage(c tele.Context, text string) error {
	return c.Send(text, &tele.SendOptions{ParseMode: tele.ModeHTML})
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
	lastMsg, ok := lastMessages[id]
	if !ok {
		return b.sendMessage(c, "Не удалось найти задачу. Убедитесь, что вы отправили сообщение перед командой.")
	}

	// Извлекаем интервал и продолжительность
	duration, _ := strconv.Atoi(matches[1])
	timeFormat := matches[2]

	err := remind.ValidateReminderDuration(timeFormat, duration)
	if err != nil {
		return b.sendMessage(c, "Неверный формат срока напоминания")
	}

	taskDuration, timeUnit := remind.CalculateReminderTime(timeFormat, duration)

	task := entity.Task{
		UserID:       id,
		Text:         lastMsg,
		CreatedAt:    time.Now(),
		Expiration:   time.Now().Add(taskDuration),
		Duration:     taskDuration,
		ReminderSent: false,
		ChatID:       c.Chat().ID,
		MsgID:        c.Message().ID,
	}

	// Отправляем ответ пользователю
	text := fmt.Sprintf("#Задача# \"%s\" принята. Напомню через <b>%d %s</b>.", task.Text, duration, timeUnit)
	_, err = b.Send(tele.ChatID(task.UserID), text, &tele.SendOptions{ParseMode: tele.ModeHTML})
	if err != nil {
		if errors.Is(err, tele.ErrChatNotFound) || errors.Is(err, tele.ErrNotStartedByUser) {
			return c.Reply(fmt.Sprintf("Для того, чтобы я запоминал твои задачи, перейди ко мне в профиль и нажми кнопку \"Запустить\".\n@%s\nЗатем вернись сюда и повтори ввод задачи и команды.", botName))
		}
		if errors.Is(err, tele.ErrBlockedByUser) {
			return c.Reply(fmt.Sprintf("Чтобы я мог писать тебе в личные сообщения, перейди по ссылке @%s и разблокируй меня.", botName))
		}
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	// Добавляем задачу с использованием сервиса
	task.ID, err = b.Services.Tasks.AddTask(ctx, task)
	if err != nil {
		return b.sendMessage(c, fmt.Sprintf("Ошибка при добавлении задачи: %s", err.Error()))
	}
	return nil
}

func (b *Bot) InitHandlers() {
	b.Bot.Handle("/start", b.Start)
	b.Bot.Handle("/help", b.Help)
	b.Bot.Handle(tele.OnText, b.HandleCtrlCommand) // Добавляем обработчик команды
	b.Bot.Handle("/tasks", b.MyTasksHandler)
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
			TaskId:     strconv.FormatInt(task.ID, 10),
			Text:       task.Text,
			CreatedAt:  task.CreatedAt,
			Expiration: task.Expiration,
			Complete:   task.ReminderSent,
		}
		text := m.String()
		return b.sendMessage(c, text)
	}
	return nil
}

func (b *Bot) Help(c tele.Context) error {
	return b.sendMessage(c, fmt.Sprintf("Чтобы сохранить задачу и создать напоминание необходимо отправить в чат свою задачу и отправить следующее сообщение в чат в формате \"<b>@%s ctrl 5d</b>\".\nГде \"5\" это интервал, а \"d\" - промежуток времени в днях. Поддерживается несколько промежутков:\nh - часы;\nd - дни;\nw - недели;\nm - месяцы.", c.Bot().Me.Username))
}
