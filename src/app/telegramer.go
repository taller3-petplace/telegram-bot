package app

import (
	"fmt"
	env "github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"os"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/db"
	"time"
)

const (
	token = "TELEGRAM_BOT_TOKEN"
)

type Telegramer struct {
	telegramBot *bot.TelegramBot
}

func NewTelegramer() (*Telegramer, error) {
	err := env.Load()
	if err != nil {
		fmt.Println("error loading environment variables")
		return nil, err
	}

	botToken := os.Getenv(token)
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

	fakeDB := db.NewFakeDB()
	telegramBot := bot.NewTelegramBot(botInstance, fakeDB)

	return &Telegramer{
		telegramBot: telegramBot,
	}, nil
}

func (t *Telegramer) Start() error {
	t.telegramBot.DefineHandlers()
	t.telegramBot.StartBot()
	return nil
}
