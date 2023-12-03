package bot

import (
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/db"
)

const (
	startEndpoint = "/start"
	infoEndpoint  = "/info"
)

// TelegramBot ToDo: add documentation
type TelegramBot struct {
	bot *tele.Bot
	db  *db.FakeDB
}

func NewTelegramBot(bot *tele.Bot, db *db.FakeDB) *TelegramBot {
	return &TelegramBot{
		bot: bot,
		db:  db,
	}
}

// DefineHandlers defines all methods that  TelegramBot can handle, is a not-blocking function
func (tb *TelegramBot) DefineHandlers() {
	// Endpoint handlers
	tb.bot.Handle(startEndpoint, tb.start)

	tb.bot.Handle(infoEndpoint, tb.info)

	// Button handlers
	tb.bot.Handle(&button.CreateAccount, tb.createAccount)

	tb.bot.Handle(&button.DontCreateAccount, tb.omitAccountCreation)
}

func (tb *TelegramBot) StartBot() {
	tb.bot.Start()
}
