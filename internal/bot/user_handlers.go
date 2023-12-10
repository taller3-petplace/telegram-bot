package bot

import (
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/utils/formatter"
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
		button.Menu.Inline(
			button.Menu.Row(button.CreateAccount),
			button.Menu.Row(button.DontCreateAccount),
		)
		message := fmt.Sprintf("You are not registered in %s, "+
			"do you want to create an account now?",
			formatter.Bold(appName),
		)
		message += fmt.Sprintf(
			"\n\n%s If you don't have an account you will not be able to perform operations with %s",
			emoji.Eyes,
			formatter.Italic(botName),
		)
		return c.Send(message, button.Menu)
	}

	welcomeMessage := template.WelcomeMessage(userInfo.GetName())
	return c.Send(welcomeMessage)
}

// createAccount returns a URL to register for appName
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
