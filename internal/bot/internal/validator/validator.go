package validator

import (
	"fmt"
	"regexp"
	"strings"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/formatter"
	"time"
)

var (
	layout     = "2006/01/02"
	validTypes = []string{
		"monkey",
		"gorilla",
		"orangutan",
		"dog",
		"poodle",
		"wolf",
		"fox",
		"raccoon",
		"cat",
		"lion",
		"tiger",
		"leopard",
		"horse",
		"zebra",
		"deer",
		"bison",
		"ox",
		"water buffalo",
		"cow",
		"pig",
		"boar",
		"ram",
		"ewe",
		"goat",
		"camel",
		"llama",
		"giraffe",
		"elephant",
		"mammoth",
		"rhinoceros",
		"hippopotamus",
		"mouse",
		"rat",
		"hamster",
		"rabbit",
		"chipmunk",
		"beaver",
		"hedgehog",
		"bat",
		"bear",
		"polar bear",
		"koala",
		"panda",
		"sloth",
		"otter",
		"skunk",
		"kangaroo",
		"badger",
		"paw",
		"turkey",
		"chicken",
		"rooster",
		"bird",
		"penguin",
		"dove",
		"eagle",
		"duck",
		"swan",
		"owl",
		"dodo",
		"feather",
		"flamingo",
		"peacock",
		"parrot",
		"frog",
		"crocodile",
		"turtle",
		"lizard",
		"snake",
		"dragon",
		"sauropod",
		"T-Rex",
		"whale",
		"dolphin",
		"seal",
		"fish",
		"blowfish",
		"shark",
		"octopus",
	}
)

// ValidatePetType returns an error if the given type is not within the valid type of pets
func ValidatePetType(petType string) error {
	petType = strings.ToLower(petType)
	if !utils.Contains(validTypes, petType) {
		return fmt.Errorf("invalid pet type: valid types are %s", strings.Join(validTypes, ", "))
	}

	return nil
}

// ValidateDateType checks if the format is year/month/day
func ValidateDateType(rawDate string) error {
	date, err := time.Parse(layout, rawDate)
	if err != nil {
		return err
	}

	diff := utils.CalculateYearsBetweenDates(date)
	if diff < 0 {
		return fmt.Errorf("error date is from the future: ")
	}

	return nil
}

// ValidateHour checks if the given hour is valid based on the next rules:
// + Hour: integer between [0, 23]
//
// + Minutes: 0 or 30
func ValidateHour(hour string) error {
	regex := regexp.MustCompile(`(^0?[0-9]|1[0-9]|2[0-3]):(00|30)`)
	if regex.MatchString(hour) {
		return nil
	}

	validHourMessage := formatter.UnorderedList([]string{"Hour: integer between [0, 23]", "Minutes: 0 or 30. E.g: 10:30, 17:00"})

	return fmt.Errorf("invalid hour for alarm: valid hour format is\n%s", validHourMessage)
}
