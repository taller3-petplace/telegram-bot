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

func (tb *TelegramBot) medicalHistory(c tele.Context) error {
	return c.Send("implement me")
}
