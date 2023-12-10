package bot

import (
	"errors"
	"fmt"
	"github.com/enescakir/emoji"
	tele "gopkg.in/telebot.v3"
	"regexp"
	"strings"
	"telegram-bot/internal/bot/internal/salchifacts"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/bot/internal/validator"
	"telegram-bot/internal/utils/formatter"
	"time"
)

const (
	// Regex group tags
	nameTag      = "Name"
	birthDateTag = "BirthDate"
	typeTag      = "Type"
)

type PetRequest struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	RegisterDate time.Time `json:"register_date"`
	BirthDate    string    `json:"birth_date"`
	OwnerID      int64     `json:"owner_id"`
}

func NewPetRequest(petData map[string]string, userID int64) PetRequest {
	date := strings.Split(petData[birthDateTag], "/")
	dateFormatted := strings.Join(date, "-")
	return PetRequest{
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
	fmt.Printf("La info de tu mascota es: %+v\n", petRequest)

	return c.Send("Pet record created correctly")
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
