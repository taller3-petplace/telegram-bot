package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"os"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/requester"
	"telegram-bot/internal/sender"
	"time"
)

const (
	tokenKey      = "TELEGRAM_BOT_TOKEN"
	senderPortKey = "SENDER_PORT"
)

type notificationSender interface {
	RegisterRoutes(r *gin.Engine)
	TriggerNotifications(c *gin.Context)
}

type App struct {
	telegramBot         *bot.TelegramBot
	notificationsSender notificationSender
}

func NewApp() (*App, error) {
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

	return &App{
		telegramBot:         telegramBot,
		notificationsSender: sender.NewNotificationSender(telegramBot),
	}, nil
}

func (a *App) RegisterRoutes(r *gin.Engine) {
	a.telegramBot.DefineHandlers()
	a.notificationsSender.RegisterRoutes(r)
}

func (a *App) Run(r *gin.Engine) error {
	errChannel := make(chan error, 1)
	go func() {
		logrus.Info("Starting bot")
		a.telegramBot.StartBot()
	}()

	port := os.Getenv(senderPortKey)
	if port == "" {
		logrus.Info("Using default port (8080) for notification sender")
		port = "8080"
	}
	go func() {
		logrus.Info("Starting notification sender")
		errChannel <- r.Run(fmt.Sprintf(":%s", port))
	}()

	err := <-errChannel
	return err
}
