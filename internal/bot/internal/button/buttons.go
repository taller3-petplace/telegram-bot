package button

import (
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
)

const (
	signInURLTemplate         = "https://web.telegram.org/a/#%d"
	createAccountEndpoint     = "create-account"
	dontCreateAccountEndpoint = "bye-dude-good-luck"
	petInfo                   = "pet-info"
	vaccines                  = "vaccines"
	medicalHistory            = "medical-history"
)

var (
	Menu              = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	PetInfoMenu       = &tele.ReplyMarkup{}
	CreateAccount     = Menu.Data("Yes", createAccountEndpoint)
	DontCreateAccount = Menu.Data("No", dontCreateAccountEndpoint)

	// PetInfo use to create different buttons for each pet of the user
	PetInfo              = Menu.Data("", petInfo)
	VaccinesButton       = Menu.Data(fmt.Sprintf("Vaccines %s", emoji.Syringe), vaccines)
	MedicalHistoryButton = Menu.Data(fmt.Sprintf("Medical history %v", emoji.OrangeBook), medicalHistory)
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
