package bot

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/requester"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/formatter"
)

const treatmentsThreshold = 5

// showVaccines shows all the vaccines that were applied to the pet. The vaccines are ordered from most recent to oldest
func (tb *TelegramBot) showVaccines(c tele.Context) error {
	params := strings.Split(c.Data(), "|")

	if len(params) != 1 {
		return c.Send(template.TryAgainMessage())
	}

	petID := params[0]
	petIDInt, err := strconv.Atoi(petID)
	if err != nil {
		fmt.Printf("invalid petID: %s\n", petID)
		return c.Send(template.TryAgainMessage())
	}

	vaccines, err := tb.requester.GetVaccines(petIDInt)

	var requestError requester.RequestError
	ok := errors.As(err, &requestError)
	if ok && (requestError.IsNotFound() || requestError.IsNoContent()) {
		return c.Send("Cannot find vaccines for selected pet")
	}

	if err != nil {
		logrus.Errorf("error fetching vaccines: petID: %s - error: %v", petID, err)
		return c.Send(template.TryAgainMessage())
	}

	message := ""
	for _, vaccine := range vaccines {
		message += fmt.Sprintf("%s\n", formatter.Bold(vaccine.Name))

		doseDates := []string{
			fmt.Sprintf("Amount of doses applied: %v", vaccine.AmountOfDoses),
			fmt.Sprintf("First Dose: %s", utils.DateToString(vaccine.FirstDose)),
			fmt.Sprintf("Last Dose: %s", utils.DateToString(vaccine.LastDose)),
		}

		message += formatter.UnorderedList(doseDates)
	}

	return c.Send(message)
}

// medicalHistory list the last 5 treatments of the pet
func (tb *TelegramBot) medicalHistory(c tele.Context) error {
	// Todo: add function to extract IDs from c.Data()
	params := strings.Split(c.Data(), "|")

	if len(params) != 1 {
		logrus.Errorf("invalid amount of params in medicalHistory: %s", params)
		return c.Send(template.TryAgainMessage())
	}

	petID := params[0]
	petIDInt, err := strconv.Atoi(petID)
	if err != nil {
		logrus.Errorf("invalid petID: %s", petID)
		return c.Send(template.TryAgainMessage())
	}

	allPetTreatments, err := tb.requester.GetTreatmentsByPetID(petIDInt)

	var requestError requester.RequestError
	ok := errors.As(err, &requestError)
	if ok && (requestError.IsNotFound() || requestError.IsNoContent()) {
		return c.Send("Cannot find treatments for selected pet")
	}

	if err != nil {
		logrus.Errorf("error fetching treatments: petID: %s - error: %v", petID, err)
		return c.Send(template.TryAgainMessage())
	}

	if len(allPetTreatments) == 0 {
		return c.Send("Your pet does not have any treatment yet")
	}

	if len(allPetTreatments) > treatmentsThreshold {
		allPetTreatments = allPetTreatments[:treatmentsThreshold]
	}

	treatmentsMenu := tb.bot.NewMarkup()
	var treatmentRows []tele.Row
	for _, treatmentData := range allPetTreatments {
		infoCut := ""
		if len(treatmentData.Comments) > 0 {
			infoCut = formatter.EllipseText(treatmentData.Comments[0].Information, 15)
		}

		buttonText := fmt.Sprintf(
			"%s: %s", treatmentData.GetName(), infoCut)

		treatmentButton := button.TreatmentSummaryButton(buttonText, treatmentData.ID)
		treatmentRows = append(treatmentRows, treatmentsMenu.Row(treatmentButton))
	}

	treatmentsMenu.Inline(treatmentRows...)

	return c.Send("Select a treatment:", treatmentsMenu)
}

// getTreatment shows all the information related with a treatment. Eg of treatment message:
// Medical appointment: 2024/01/08
// Next Turn: 2024/02/20 or -
// Date End: 2024/02/10 or -
// Comments:
//   - 2023/12/18 by Lasso: nada es igual, nada es igual sin tus ojos marrones
//   - 2023/10/05 by Arjona: tu reputacion son las primeras seis letras de esa palabra
func (tb *TelegramBot) getTreatment(c tele.Context) error {
	params := strings.Split(c.Data(), "|")

	if len(params) != 1 {
		logrus.Errorf("receive invalid amount of data in getTreatment: %v", params)
		return c.Send(template.TryAgainMessage())
	}

	treatmentID := params[0]
	treatment, err := tb.requester.GetTreatment(treatmentID)

	var requestError requester.RequestError
	ok := errors.As(err, &requestError)
	if ok && (requestError.IsNotFound() || requestError.IsNoContent()) {
		return c.Send("Cannot find info about selected treatment")
	}

	if err != nil {
		logrus.Errorf("error fetching treatment: treatmentID: %s - error: %v\n", treatmentID, err)
		return c.Send(template.TryAgainMessage())
	}

	dateEnd := "-"
	if treatment.DateEnd != nil {
		dateEnd = utils.DateToString(*treatment.DateEnd)
	}

	nextTurn := "-"
	if treatment.NextTurn != nil {
		nextTurn = utils.DateToString(*treatment.NextTurn)
	}

	message := fmt.Sprintf(
		"%s\n\nNext Turn: %s \nDate End: %s \nComments:\n",
		formatter.Bold(treatment.GetName()),
		nextTurn,
		dateEnd,
	)

	var commentMessages []string
	for _, comment := range treatment.Comments {
		commentMessages = append(commentMessages, comment.GetCommentMessage())
	}

	message += formatter.UnorderedList(commentMessages)
	return c.Send(message)
}
