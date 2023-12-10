package bot

import (
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/db"
)

const (
	appName     = "Pet Place"
	botName     = "Ringot"
	botUsername = "@pet_place_bot"

	// Endpoints
	startEndpoint       = "/start"
	helpEndpoint        = "/help"
	createPetEndpoint   = "/createPet"
	registerPetEndpoint = "/addPetRecord"
	salchiFactEndpoint  = "/salchiFact"
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
	tb.bot.Handle(helpEndpoint, tb.help)

	tb.bot.Handle(startEndpoint, tb.start)

	tb.bot.Handle(createPetEndpoint, tb.createPet)

	tb.bot.Handle(salchiFactEndpoint, tb.getSalchiFact)

	// Button handlers
	tb.bot.Handle(&button.CreateAccount, tb.createAccount)

	tb.bot.Handle(&button.DontCreateAccount, tb.omitAccountCreation)

	// Action handlers
	tb.bot.Handle(tele.OnText, tb.textHandler)

	tb.bot.Handle(tele.OnEdited, tb.editMessageHandler)
}

func (tb *TelegramBot) StartBot() {
	tb.bot.Start()
}
