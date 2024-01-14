package bot

import (
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/requester"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/formatter"
)

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
	if ok && requestError.IsNotFound() || requestError.IsNoContent() {
		return c.Send("Cannot find vaccines for selected pet")
	}

	if err != nil {
		fmt.Printf("error fetching vaccines: petID: %s - error: %v\n", petID, err)
		return c.Send(template.TryAgainMessage())
	}

	message := ""
	for _, vaccine := range vaccines {
		message += fmt.Sprintf("%s\n", formatter.Bold(vaccine.Name))

		nextDose := "-"
		if vaccine.NextDose != nil {
			nextDose = utils.DateToString(*vaccine.NextDose)
		}

		doseDates := []string{
			utils.DateToString(vaccine.FirstDose),
			utils.DateToString(vaccine.LastDose),
			nextDose,
		}

		message += formatter.UnorderedList(doseDates)
	}

	return c.Send(message)
}

// medicalHistory list the last 5 treatments of the pet
func (tb *TelegramBot) medicalHistory(c tele.Context) error {
	return c.Send("implement me")
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
		return c.Send(template.TryAgainMessage())
	}

	treatmentID := params[0]
	treatmentIDInt, err := strconv.Atoi(treatmentID)
	if err != nil {
		fmt.Printf("invalid petID: %s\n", treatmentID)
		return c.Send(template.TryAgainMessage())
	}

	treatment, err := tb.requester.GetTreatment(treatmentIDInt)

	var requestError requester.RequestError
	ok := errors.As(err, &requestError)
	if ok && requestError.IsNotFound() || requestError.IsNoContent() {
		return c.Send("Cannot find info about selected treatment")
	}

	if err != nil {
		fmt.Printf("error fetching treatment: treatmentID: %s - error: %v\n", treatmentID, err)
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
		"%s\nNext Turn: %s \nDate End: %s \nComments:",
		treatment.GetName(),
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
