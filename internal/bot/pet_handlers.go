package bot

import (
	"errors"
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
	"regexp"
	"strings"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/bot/internal/salchifacts"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/bot/internal/validator"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/requester"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/formatter"
	"time"
)

const (
	// Regex group tags
	nameTag      = "Name"
	birthDateTag = "BirthDate"
	typeTag      = "Type"
	hoursInAYear = 365 * 24
)

func NewPetRequest(petData map[string]string, userID int64) domain.PetRequest {
	date := strings.Split(petData[birthDateTag], "/")
	dateFormatted := strings.Join(date, "-")
	return domain.PetRequest{
		Name:         formatter.Capitalize(petData[nameTag]),
		Type:         strings.ToLower(petData[typeTag]),
		BirthDate:    dateFormatted,
		OwnerID:      userID,
		RegisterDate: time.Now(),
	}
}

var tryAgainMessage = "Try again editing the form message or execute /createPet to start again"

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
	senderInfo := c.Sender()
	if senderInfo == nil {
		return fmt.Errorf("error sender info not found")
	}

	petData, err := extractPetData(c.Message().Text)
	if err != nil && errors.Is(err, errInvalidPetForm) {
		return c.Send(fmt.Sprintf("%v Invalid form, you don't have to modify the structure, only the field values. %s",
			emoji.PoliceCarLight,
			tryAgainMessage,
		))
	}

	if err != nil && errors.Is(err, errMissingPetField) {
		return c.Send("%v %v. Try again editing the form message and adding the deleted field or execute /cretePet", emoji.PoliceCarLight, err)
	}

	if err := validator.ValidatePetType(petData[typeTag]); err != nil {
		return c.Send(fmt.Sprintf("%v. %s", err, tryAgainMessage))
	}

	if err := validator.ValidateDateType(petData[birthDateTag]); err != nil {
		return c.Send(fmt.Sprintf("Invalid birth date: format must be year/month/day. %s", tryAgainMessage))
	}

	if len(petData[nameTag]) == 0 {
		return c.Send(fmt.Sprintf("The most important thing is missing, the name of your pet! %s", tryAgainMessage))
	}

	petRequest := NewPetRequest(petData, senderInfo.ID)

	err = tb.requester.RegisterPet(petRequest)
	if err != nil {
		fmt.Printf("hubo un error creando la mascota: %v", err)
		return c.Send("error creando la mascota al hacer la request")
	}

	return c.Send("Pet record created correctly")
}

// getPets search for the owner's pets based on telegram ID
func (tb *TelegramBot) getPets(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		return fmt.Errorf("error sender info not found")
	}

	petsData, err := tb.requester.GetPetsByOwnerID(senderInfo.ID)

	var requestError requester.RequestError
	ok := errors.As(err, &requestError)
	if ok && requestError.IsNotFound() {
		return c.Send("You don't have any pet registered yet")
	}

	if err != nil {
		fmt.Printf("\n el error: %v", err)
		return c.Send("error searching your pets")
	}

	petsMenu := tb.bot.NewMarkup()

	var petRows []tele.Row
	for _, petData := range petsData {
		petEmoji := utils.GetEmojiForPetType(petData.Type)
		buttonText := fmt.Sprintf("%s %v", petData.Name, petEmoji)

		petButton := petsMenu.Data(buttonText, button.PetInfo.Unique, fmt.Sprintf("%v", petData.ID))
		petRows = append(petRows, petsMenu.Row(petButton))
	}

	petsMenu.Inline(petRows...)

	return c.Send("Select a pet", petsMenu)
}

func (tb *TelegramBot) getPetInfo(c tele.Context) error {
	petData := strings.Split(c.Data(), "|")

	// petData.Name, fmt.Sprintf("%v", petData.ID), petData.Type, age

	message := fmt.Sprintf("%s \n\n", petData[0])
	petInfoItems := []string{
		fmt.Sprintf("Age: %v", petData[3]),
		fmt.Sprintf("Type: %s", petData[2]),
	}

	message += formatter.UnorderedList(petInfoItems)

	petInfoMenu := tb.bot.NewMarkup()

	petInfoMenu.Inline(
		petInfoMenu.Row(button.MedicalHistoryButton),
		petInfoMenu.Row(button.VaccinesButton),
	)

	return c.Send(message, petInfoMenu)
}

func (tb *TelegramBot) showVaccines(c tele.Context) error {
	return nil
}

func (tb *TelegramBot) medicalHistory(q *tele.Query) error {
	return tb.bot.Answer(q, &tele.QueryResponse{})
}

// getSalchiFact returns a random fact about perros salchichas
func (tb *TelegramBot) getSalchiFact(c tele.Context) error {
	fact := salchifacts.GetFact()
	return c.Send(fact)
}

// extractPetData extracts pet data from the given message. Does not validate the fields, it only ensures that they are all present
func extractPetData(petDataRaw string, fields ...string) (map[string]string, error) {
	regex := regexp.MustCompile(`Name:\s?(?P<Name>[^\n]*)\s+Birth Date:\s?(?P<BirthDate>[^\n]*)\s+Type:\s?(?P<Type>[^\n]*)`)
	match := regex.FindStringSubmatch(petDataRaw)
	if match == nil {
		return nil, fmt.Errorf("%w", errInvalidPetForm)
	}

	// groupName are capture from the regex expression
	petData := make(map[string]string)
	for idx, groupName := range regex.SubexpNames() {
		if idx != 0 && groupName != "" {
			petData[groupName] = match[idx]
		}
	}

	for _, field := range fields {
		if _, found := petData[field]; !found {
			return nil, fmt.Errorf("%w: %s", errMissingPetField, field)
		}
	}

	return petData, nil
}
