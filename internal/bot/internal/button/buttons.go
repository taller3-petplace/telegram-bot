package button

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

const (
	signInURLTemplate = "https://web.telegram.org/a/#%d"
)

var (
	Menu              = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	CreateAccount     = Menu.Text("Yes")
	DontCreateAccount = Menu.Text("No")
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
