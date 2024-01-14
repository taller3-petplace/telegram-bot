package utils

import (
	"fmt"
	"github.com/enescakir/emoji"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

const hoursInAYear = 365 * 24

var animalEmojisMap = map[string]emoji.Emoji{
	"monkey":        emoji.Monkey,
	"gorilla":       emoji.Gorilla,
	"orangutan":     emoji.Orangutan,
	"dog":           emoji.DogFace,
	"poodle":        emoji.Poodle,
	"wolf":          emoji.Wolf,
	"fox":           emoji.Fox,
	"raccoon":       emoji.Raccoon,
	"cat":           emoji.CatFace,
	"lion":          emoji.Lion,
	"tiger":         emoji.Tiger,
	"leopard":       emoji.Leopard,
	"horse":         emoji.Horse,
	"zebra":         emoji.Zebra,
	"deer":          emoji.Deer,
	"bison":         emoji.Bison,
	"ox":            emoji.Ox,
	"water buffalo": emoji.WaterBuffalo,
	"cow":           emoji.Cow,
	"pig":           emoji.Pig,
	"boar":          emoji.Boar,
	"ram":           emoji.Ram,
	"ewe":           emoji.Ewe,
	"goat":          emoji.Goat,
	"camel":         emoji.Camel,
	"llama":         emoji.Llama,
	"giraffe":       emoji.Giraffe,
	"elephant":      emoji.Elephant,
	"mammoth":       emoji.Mammoth,
	"rhinoceros":    emoji.Rhinoceros,
	"hippopotamus":  emoji.Hippopotamus,
	"mouse":         emoji.Mouse,
	"rat":           emoji.Rat,
	"hamster":       emoji.Hamster,
	"rabbit":        emoji.Rabbit,
	"chipmunk":      emoji.Chipmunk,
	"beaver":        emoji.Beaver,
	"hedgehog":      emoji.Hedgehog,
	"bat":           emoji.Bat,
	"bear":          emoji.Bear,
	"polar bear":    emoji.PolarBear,
	"koala":         emoji.Koala,
	"panda":         emoji.Panda,
	"sloth":         emoji.Sloth,
	"otter":         emoji.Otter,
	"skunk":         emoji.Skunk,
	"kangaroo":      emoji.Kangaroo,
	"badger":        emoji.Badger,
	"paw":           emoji.PawPrints,
	"turkey":        emoji.Turkey,
	"chicken":       emoji.Chicken,
	"rooster":       emoji.Rooster,
	"bird":          emoji.Bird,
	"penguin":       emoji.Penguin,
	"dove":          emoji.Dove,
	"eagle":         emoji.Eagle,
	"duck":          emoji.Duck,
	"swan":          emoji.Swan,
	"owl":           emoji.Owl,
	"dodo":          emoji.Dodo,
	"feather":       emoji.Feather,
	"flamingo":      emoji.Flamingo,
	"peacock":       emoji.Peacock,
	"parrot":        emoji.Parrot,
	"frog":          emoji.Frog,
	"crocodile":     emoji.Crocodile,
	"turtle":        emoji.Turtle,
	"lizard":        emoji.Lizard,
	"snake":         emoji.Snake,
	"dragon":        emoji.Dragon,
	"sauropod":      emoji.Sauropod,
	"t-rex":         emoji.TRex,
	"whale":         emoji.Whale,
	"dolphin":       emoji.Dolphin,
	"seal":          emoji.Seal,
	"fish":          emoji.Fish,
	"blowfish":      emoji.Blowfish,
	"shark":         emoji.Shark,
	"octopus":       emoji.Octopus,
}

func ReadFileWithPath(configFilePath string, suffixToRemove string) ([]byte, error) {
	_, file, _, _ := runtime.Caller(1)
	filePath := strings.TrimSuffix(file, suffixToRemove)
	filePath += configFilePath

	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	configFileBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return configFileBytes, nil
}

func Contains[T comparable](elements []T, target T) bool {
	for idx := range elements {
		if elements[idx] == target {
			return true
		}
	}

	return false
}

func GetEmojiForPetType(petType string) emoji.Emoji {
	petType = strings.ToLower(petType)
	return animalEmojisMap[petType]
}

// CalculateYearsBetweenDates calculates the amount of years between the given date and the current one
func CalculateYearsBetweenDates(date time.Time) int {
	diff := time.Now().Sub(date)

	amountOfYears := diff.Hours() / hoursInAYear
	return int(amountOfYears)
}

type sorter interface {
	GetDate() time.Time
}

// SortElementsByDate sorts from newest dates to oldest ones
func SortElementsByDate[T sorter](elements []T) {
	sort.Slice(elements, func(i, j int) bool {
		date1 := elements[i].GetDate()
		date2 := elements[j].GetDate()

		return date1.After(date2)
	})
}

// DateToString transforms the input in a string with format yyyy-mm-dd
func DateToString(date time.Time) string {
	return date.Format(time.DateOnly)
}
