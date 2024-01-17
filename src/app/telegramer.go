package app

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"os"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/requester"
	"time"
)

const (
	tokenKey = "TELEGRAM_BOT_TOKEN"
)

type Telegramer struct {
	telegramBot *bot.TelegramBot
}

func NewTelegramer() (*Telegramer, error) {
	botToken := os.Getenv(tokenKey)
	if botToken == "" {
		return nil, fmt.Errorf("bot token is missing")
	}

	botSettings := tele.Settings{
		Token:     botToken,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeMarkdown,
	}

	botInstance, err := tele.NewBot(botSettings)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	serviceRequester, err := requester.NewRequester(client)
	if err != nil {
		return nil, err
	}

	telegramBot := bot.NewTelegramBot(botInstance, serviceRequester)

	return &Telegramer{
		telegramBot: telegramBot,
	}, nil
}

func (t *Telegramer) Start() error {
	t.telegramBot.DefineHandlers()
	t.telegramBot.StartBot()
	return nil
}
