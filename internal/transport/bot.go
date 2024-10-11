package transport

import (
	"fmt"
	"time"

	"tasktrackerbot/internal/service"

	tele "gopkg.in/telebot.v3"
)

type BotService struct {
	*tele.Bot
	Services service.Services
}

func (b BotService) Start() {
	fmt.Println("Запускаю")
	b.Bot.Start()
}

func NewBotService(token string, services service.Services) BotService {
	prefs := tele.Settings{
		Token:  token,
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
