package bot

import (
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/requester"
)

const (
	appName     = "Pet Place"
	botName     = "Ringot"
	botUsername = "@pet_place_bot"

	// Endpoints
	startEndpoint                = "/start"
	helpEndpoint                 = "/help"
	createPetEndpoint            = "/createPet"
	getPets                      = "/getPets"
	registerPetEndpoint          = "/addPetRecord"
	salchiFactEndpoint           = "/salchiFact"
	setNotificationEndpoint      = "/setNotification"
	registerNotificationEndpoint = "/notification"
	getVetsEndpoint              = "/getVets"
)

// TelegramBot handles requests from telegram. Is in charge to interact with different services
// in order to give a response to the request of the user.
//
// What this bot can do is defined in DefineHandlers
type TelegramBot struct {
	bot       *tele.Bot
	usersDB   map[int64]bool
	requester *requester.Requester
}

func NewTelegramBot(bot *tele.Bot, requester *requester.Requester) *TelegramBot {
	usersDB := make(map[int64]bool)
	return &TelegramBot{
		bot:       bot,
		requester: requester,
		usersDB:   usersDB,
	}
}

// DefineHandlers defines all methods that  TelegramBot can handle, is a not-blocking function
func (tb *TelegramBot) DefineHandlers() {
	// Endpoints handlers
	tb.bot.Handle(helpEndpoint, tb.help)

	tb.bot.Handle(startEndpoint, tb.start)

	tb.bot.Handle(createPetEndpoint, tb.createPet)

	tb.bot.Handle(getPets, tb.getPets)

	tb.bot.Handle(salchiFactEndpoint, tb.getSalchiFact)

	tb.bot.Handle(getVetsEndpoint, tb.getVets)

	tb.bot.Handle(setNotificationEndpoint, tb.setAlarm)

	// Button handlers
	tb.bot.Handle(&button.CreateAccount, tb.createAccount)

	tb.bot.Handle(&button.DontCreateAccount, tb.omitAccountCreation)

	tb.bot.Handle(&button.PetInfo, tb.getPetInfo)

	tb.bot.Handle(&button.MedicalHistory, tb.medicalHistory)

	tb.bot.Handle(&button.Vaccines, tb.showVaccines)

	tb.bot.Handle(&button.Treatment, tb.getTreatment)

	// Action handlers
	tb.bot.Handle(tele.OnText, tb.textHandler)

	tb.bot.Handle(tele.OnEdited, tb.editMessageHandler)

	tb.bot.Handle(tele.OnLocation, tb.searchVets)
}

func (tb *TelegramBot) StartBot() {
	tb.bot.Start()
}

func (tb *TelegramBot) SendNotification(telegramID int64, messageBody string) error {
	chat, err := tb.bot.ChatByID(telegramID)
	if err != nil {
		return fmt.Errorf("error fetching chat of user %d: %v", telegramID, err)
	}
	message := fmt.Sprintf("Scheduled notification %s\n%s", emoji.AlarmClock, messageBody)
	_, err = tb.bot.Send(chat, message)
	if err != nil {
		return err
	}

	return nil
}
