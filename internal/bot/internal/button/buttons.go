package button

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

const (
	signInURLTemplate         = "https://web.telegram.org/a/#%d"
	createAccountEndpoint     = "create-account"
	dontCreateAccountEndpoint = "bye-dude-good-luck"
)

var (
	Menu              = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	CreateAccount     = Menu.Data("Yes", createAccountEndpoint)
	DontCreateAccount = Menu.Data("No", dontCreateAccountEndpoint)
)

func SignUpButton(telegramID int64) *tele.ReplyMarkup {
	signUpButton := &tele.ReplyMarkup{}

	url := fmt.Sprintf(signInURLTemplate, telegramID)
	buttonURL := signUpButton.URL("Sign Up", url)

	signUpButton.Inline(
		signUpButton.Row(buttonURL),
	)

	return signUpButton
}
