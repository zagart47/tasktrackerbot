package transport

import (
	"time"

	"tasktrackerbot/config"
	"tasktrackerbot/internal/service"

	tele "gopkg.in/telebot.v3"
)

type BotService struct {
	*tele.Bot
	Services service.Services
}

func NewBotService(services service.Services) BotService {
	prefs := tele.Settings{
		Token:  config.Configs.Bot.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(prefs)
	if err != nil {
		panic(err)
	}

	return BotService{
		Bot:      b,
		Services: services,
	}
}
