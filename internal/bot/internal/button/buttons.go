package button

import (
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
)

const (
	signInURLTemplate         = "https://lnt.digital/#/sign-up?telegram_id=%d"
	createAccountEndpoint     = "create-account"
	dontCreateAccountEndpoint = "bye-dude-good-luck"
	petInfoEndpoint           = "pet-info"
	vaccinesEndpoint          = "vaccines"
	medicalHistoryEndpoint    = "medical-history"
	setAlarmEndpoint          = "set-alarm"
	treatmentInfoEndpoint     = "treatment-info"
)

var (
	Menu              = &tele.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	CreateAccount     = Menu.Data("Yes", createAccountEndpoint)
	DontCreateAccount = Menu.Data("No", dontCreateAccountEndpoint)

	// PetInfo use to create different buttons for each pet of the user
	PetInfo        = Menu.Data("", petInfoEndpoint)
	Vaccines       = Menu.Data(fmt.Sprintf("Vaccines %s", emoji.Syringe), vaccinesEndpoint)
	MedicalHistory = Menu.Data(fmt.Sprintf("Medical history %v", emoji.OrangeBook), medicalHistoryEndpoint)
	Treatment      = Menu.Data("", treatmentInfoEndpoint)
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

func VaccinesButton(petID string) tele.Btn {
	markup := &tele.ReplyMarkup{}
	return markup.Data(Vaccines.Text, Vaccines.Unique, petID)
}

func MedicalHistoryButton(petID string) tele.Btn {
	markup := &tele.ReplyMarkup{}
	return markup.Data(MedicalHistory.Text, MedicalHistory.Unique, petID)
}

func TreatmentSummaryButton(treatmentSummary string, treatmentID string) tele.Btn {
	markup := &tele.ReplyMarkup{}
	return markup.Data(treatmentSummary, Treatment.Unique, treatmentID)
}
