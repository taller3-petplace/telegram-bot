package validator

import (
	"fmt"
	"strings"
	"telegram-bot/internal/utils"
	"time"
)

// ToDo: define some validator functions

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
func ValidateDateType(date string) error {
	_, err := time.Parse(layout, date)
	return err
}
