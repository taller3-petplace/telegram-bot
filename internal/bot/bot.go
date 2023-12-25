package bot

import (
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/db"
	"telegram-bot/internal/requester"
)

const (
	appName     = "Pet Place"
	botName     = "Ringot"
	botUsername = "@pet_place_bot"

	// Endpoints
	startEndpoint       = "/start"
	helpEndpoint        = "/help"
	createPetEndpoint   = "/createPet"
	getPets             = "/getPets"
	registerPetEndpoint = "/addPetRecord"
	salchiFactEndpoint  = "/salchiFact"
)

// TelegramBot ToDo: add documentation
type TelegramBot struct {
	bot       *tele.Bot
	db        *db.FakeDB
	requester *requester.Requester
}

func NewTelegramBot(bot *tele.Bot, db *db.FakeDB, requester *requester.Requester) *TelegramBot {
	return &TelegramBot{
		bot:       bot,
		db:        db,
		requester: requester,
	}
}

// DefineHandlers defines all methods that  TelegramBot can handle, is a not-blocking function
func (tb *TelegramBot) DefineHandlers() {
	// Endpoint handlers
	tb.bot.Handle(helpEndpoint, tb.help)

	tb.bot.Handle(startEndpoint, tb.start)

	tb.bot.Handle(createPetEndpoint, tb.createPet)

	tb.bot.Handle(getPets, tb.getPets)

	tb.bot.Handle(salchiFactEndpoint, tb.getSalchiFact)

	// Button handlers
	tb.bot.Handle(&button.CreateAccount, tb.createAccount)

	tb.bot.Handle(&button.DontCreateAccount, tb.omitAccountCreation)

	tb.bot.Handle(&button.PetInfo, tb.getPetInfo)

	//tb.bot.Handle(&button.MedicalHistoryButton, tb.medicalHistory)

	tb.bot.Handle(&button.VaccinesButton, tb.showVaccines)

	// Action handlers
	tb.bot.Handle(tele.OnText, tb.textHandler)

	tb.bot.Handle(tele.OnEdited, tb.editMessageHandler)

	//tb.bot.Handle(tele.OnQuery, tb.medicalHistory)
}

func (tb *TelegramBot) StartBot() {
	tb.bot.Start()
}
