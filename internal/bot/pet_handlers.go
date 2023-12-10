package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"telegram-bot/internal/bot/internal/salchifact"
	"telegram-bot/internal/bot/internal/template"
)

// createPet sends a form to the user that will contain the data for the new pet
func (tb *TelegramBot) createPet(c tele.Context) error {
	petFormMenu := tb.bot.NewMarkup()
	helpButton := petFormMenu.Text("Click here to display the form")

	petForm := fmt.Sprintf("%s\n\n", registerPetEndpoint)
	petForm += template.RegisterPet()

	helpButton.InlineQueryChat = petForm

	petFormMenu.Inline(
		petFormMenu.Row(helpButton),
	)

	return c.Send("Please, enter your pet info", petFormMenu)
}

// createPetRecord creates a new record for a pet
func (tb *TelegramBot) createPetRecord(c tele.Context) error {
	// ToDo: parse info
	return c.Send("Pet record created correctly")
}

// getSalchiFact returns a random fact about perros salchichas
func (tb *TelegramBot) getSalchiFact(c tele.Context) error {
	fact := salchifact.GetFact()
	return c.Send(fact)
}
