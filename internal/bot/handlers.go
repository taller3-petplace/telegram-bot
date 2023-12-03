package bot

import (
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/utils/formatter"
)

const (
	petPlace = "Pet Place"
	botName  = "Ringot"
)

// start this endpoint has two possible flows
// 1. If the user is registered, awaits for other commands
// 2. If is not registered, gives to the user the options to create an account
func (tb *TelegramBot) start(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		return fmt.Errorf("error user info not found")
	}

	userInfo, found := tb.db.GetUser(senderInfo.ID)
	if !found {
		button.Menu.Reply(
			button.Menu.Row(button.CreateAccount),
			button.Menu.Row(button.DontCreateAccount),
		)
		message := fmt.Sprintf("You are not registered in %s, "+
			"do you want to create an account now?",
			formatter.Bold(petPlace),
		)
		message += fmt.Sprintf(
			"\n\n%s If you don't have an account you will not be able to perform operations with %s",
			emoji.Eyes,
			formatter.Italic(botName),
		)
		return c.Send(message, button.Menu)
	}

	infoMessage := getInfoMessage(userInfo.GetName())
	return c.Send(infoMessage)
}

func (tb *TelegramBot) info(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		return fmt.Errorf("error user info not found")
	}

	url := "https://www.youtube.com/watch?v=RWIJExat-lQ&list=RDRWIJExat-lQ&start_radio=1&ab_channel=Estoes%C2%A1FA%21"

	message := fmt.Sprintf(
		"Dale cachorro, escuchate %s %s %s",
		formatter.Link("esta", url),
		senderInfo.FirstName,
		emoji.MusicalNotes,
	)
	return c.Send(message)
}

// createAccount returns a URL to register for petPlace
func (tb *TelegramBot) createAccount(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		return fmt.Errorf("%w", errUserInfoNotFound)
	}

	tb.db.AddUser(*senderInfo)

	signUpButton := button.SignUpButton(senderInfo.ID)
	message := fmt.Sprintf("Click below to sign up %s", emoji.BackhandIndexPointingDown)
	err := c.Send(message, signUpButton)
	if err != nil {
		return fmt.Errorf("%w: %w", errSendingSignUpLink, err)
	}

	afterCreationMessage := fmt.Sprintf("After creating the account perform /start again %s", emoji.GrinningCatWithSmilingEyes)
	return c.Send(afterCreationMessage)
}

// omitAccountCreation byd dude, good luck
func (tb *TelegramBot) omitAccountCreation(c tele.Context) error {
	p := &tele.Photo{File: tele.FromURL("https://pbs.twimg.com/media/FRxJVLYXwAAlGPk?format=jpg&name=small")}
	_, err := p.Send(tb.bot, c.Recipient(), nil)
	return err
}

// getInfoMessage returns a message with all the information about Ringot and Pet Place
func getInfoMessage(userName string) string {
	message := fmt.Sprintf(
		"Welcome to %s, %s!. I'm %s and I'll help you to perform different operations from Telegram %s. My features are:\n\n",
		petPlace,
		userName,
		botName,
		emoji.Airplane,
	)

	urlPerroSalchicha := "https://www.youtube.com/watch?v=IQ9kDtbwoaw"
	hyperlink := formatter.Link("perro salchicha gordo bachicha", urlPerroSalchicha)

	features := []string{
		fmt.Sprintf("/start: action to start a conversation with me %s%s", emoji.DogFace, emoji.Robot),
		fmt.Sprintf("/createPet: creates a register for your pet on-demand %s", emoji.Notebook),
		fmt.Sprintf("/getPets: looks for information about your pets %s %s %s %s ", emoji.DogFace, emoji.Cat, emoji.Crocodile, emoji.Otter),
		fmt.Sprintf("/setAlarm: sets an alarm whenever you want in your timezone %s", emoji.AlarmClock),
		fmt.Sprintf("/getVets: search veterinaries %s near your location", emoji.Hospital),
		fmt.Sprintf(
			"/salchiFact: we all love '%s', so what's better that a random fact about salchichas? %s %s #SalchiData\n",
			hyperlink,
			emoji.HotDog,
			emoji.DogFace,
		),
	}

	featuresList := formatter.UnorderedList(features)
	return message + featuresList
}
