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
	return b.sendMessage(c, "Если ты видишь это сообщение, значит всё идет правильно. Теперь я буду запоминать твои задачи и напоминать тебе о них здесь.", b.getMainMenuKeyboard())
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
	lastMsg, ok := lastMessages[id]
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
		Text:         lastMsg,
		CreatedAt:    time.Now(),
		Expiration:   time.Now().Add(taskDuration),
		Duration:     taskDuration,
		ReminderSent: false,
		ChatID:       c.Chat().ID,
		MsgID:        c.Message().ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Configs.Timeout)
	defer cancel()
	// Добавляем задачу с использованием сервиса
	task.ID, err = b.Services.Tasks.AddTask(ctx, task)
	if err != nil {
		return c.Reply(fmt.Sprintf("Ошибка при добавлении задачи: %s", err.Error()))
	}

	// Отправляем ответ пользователю
	text := fmt.Sprintf("#Задача# \"%s\" принята. Присвоил ей номер - <b>%v</b>. Напомню через <b>%d %s</b>.", task.Text, task.ID, duration, timeUnit)
	botURL := fmt.Sprintf("t.me/%s", botName)
	_, err = b.Send(tele.ChatID(task.UserID), text, &tele.SendOptions{ParseMode: tele.ModeHTML})
	if err != nil {
		if errors.Is(err, tele.ErrChatNotFound) || errors.Is(err, tele.ErrNotStartedByUser) {
			return c.Reply(fmt.Sprintf("Привет! Я бот для управления задачами. Если ты будешь отправлять мне задачи в чате, я буду запоминать их и напоминать тебе о них. Для того чтобы я начал запоминать твои задачи, перейди ко мне в профиль и нажми кнопку \"Старт\".\n%s", botURL))
		}
		if errors.Is(err, tele.ErrBlockedByUser) {
			return c.Reply(fmt.Sprintf("Чтобы я мог писать тебе в личные сообщения, перейди по ссылке %s и разблокируй меня.", botURL))
		}
		return err
	}
	return nil
}

func (b *Bot) InitHandlers() {
	b.Bot.Handle("/start", b.Start)
	b.Bot.Handle("/help", b.Help)
	b.Bot.Handle(tele.OnText, b.HandleCtrlCommand) // Добавляем обработчик команды
	b.Bot.Handle(&tele.ReplyButton{Text: "Мои задачи"}, b.MyTasksHandler)
	b.Bot.Handle(&tele.InlineButton{Unique: "start_messaging"}, b.StartMessaging)
	go b.StartTasksSending()
}

func (b *Bot) StartMessaging(c tele.Context) error {
	_, err := b.Send(&tele.User{ID: b.Me.ID}, "/start")
	if err != nil {
		return err
	}
	return nil
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
		b.sendMessage(c, text, b.getMainMenuKeyboard())
	}
	return nil
}

func (b *Bot) Help(c tele.Context) error {
	return c.Send(fmt.Sprintf("Чтобы сохранить задачу и создать напоминание необходимо написать мне в чате сообщение формата \"@%s ctrl 5d\".\nГде \"5\" это интервал, а \"d\" - промежуток времени в днях. Поддерживается несколько промежутков:\nh - часы;\nd - дни;\nw - недели;\nm - месяцы.", c.Bot().Me.Username))
}
