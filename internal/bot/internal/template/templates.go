package template

import (
	"fmt"
	"github.com/enescakir/emoji"
	"telegram-bot/internal/utils/formatter"
)

// RegisterPet contains the data for the creation of a pet register
func RegisterPet() string {
	form := "Name: yourPetName\n\n"
	form += "Birth Date: yyyy/mm/dd\n\n"
	form += "Type: E.g: cat, dog, otter, etc"

	return form
}

// NotificationForm contains the data to set an alarm for a given period or even indeterminately
func NotificationForm() string {
	form := "Message: message to be sent\n\n"
	form += "Hours: hh1:mm1, hh2:mm2\n\n"
	form += "Start Date: yyyy/mm/dd\n\n"
	form += "End Date: yyyy/mm/dd or N/A"

	return form
}

func WelcomeMessage(userName string) string {
	message := fmt.Sprintf(
		"Welcome to Pet Place, %s! I'm Ringot and I'll help you to perform different operations from Telegram %s. My features are:\n\n",
		userName,
		emoji.SmallAirplane,
	)

	return message + Commands()
}

// Commands gives a list with all commands that bot supports
func Commands() string {
	urlPerroSalchicha := "https://www.youtube.com/watch?v=IQ9kDtbwoaw"
	hyperlink := formatter.Link("perro salchicha gordo bachicha", urlPerroSalchicha)

	commands := []string{
		fmt.Sprintf("/start: action to start a conversation with me %s%s", emoji.DogFace, emoji.Robot),
		fmt.Sprintf("/help: gives information about the commands that the bot supports"),
		fmt.Sprintf("/createPet: creates a register for your pet on-demand %s", emoji.Notebook),
		fmt.Sprintf("/getPets: looks for information about your pets %s %s %s %s ", emoji.DogFace, emoji.CatFace, emoji.Crocodile, emoji.Otter),
		fmt.Sprintf("/setNotification: sets an alarm whenever you want in your timezone %s", emoji.AlarmClock),
		fmt.Sprintf("/getVets: search vets %s near your location", emoji.Hospital),
		fmt.Sprintf(
			"/salchiFact: we all love '%s', so what's better that a random fact about salchichas? %s %s #SalchiData\n",
			hyperlink,
			emoji.HotDog,
			emoji.DogFace,
		),
	}

	return formatter.UnorderedList(commands)
}

func TryAgainMessage() string {
	return "Something went wrong, please try again"
}
