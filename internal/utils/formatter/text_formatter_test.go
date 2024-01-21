package formatter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBold(t *testing.T) {
	text := "hola que tal tu como estas? dime si eres feliz"
	expectedResult := fmt.Sprintf("**%s**", text)
	boldText := Bold(text)
	assert.Equal(t, expectedResult, boldText)
}

func TestItalic(t *testing.T) {
	text := "hola que tal tu como estas? dime si eres feliz"
	expectedResult := fmt.Sprintf("_%s_", text)
	italicText := Italic(text)
	assert.Equal(t, expectedResult, italicText)
}

func TestLink(t *testing.T) {
	url := "https://music.youtube.com/watch?v=42cTngAoXw4&si=q0V81YHMkrNWlrNP"
	text := "temazo"
	expectedResult := fmt.Sprintf("[%s](%s)", text, url)
	textWithLink := Link(text, url)
	assert.Equal(t, expectedResult, textWithLink)
}

func TestCapitalize(t *testing.T) {
	text := "te estas portando mal, seras castigada"
	capitalizeText := Capitalize(text)
	assert.Equal(t, "Te estas portando mal, seras castigada", capitalizeText)
}

func TestUnderlineText(t *testing.T) {
	text := "hola que tal tu como estas? dime si eres feliz"
	expectedResult := fmt.Sprintf("__%s__", text)
	italicText := UnderlineText(text)
	assert.Equal(t, expectedResult, italicText)
}

func TestSpoilerText(t *testing.T) {
	text := "Olvídala. No es fácil para mí, por eso quiero hablarle, " +
		"si es preciso rogarle que regrese a mi vida"
	expectedResult := fmt.Sprintf("||%s||", text)
	italicText := SpoilerText(text)
	assert.Equal(t, expectedResult, italicText)
}

func TestOrderedList(t *testing.T) {
	items := []string{
		"Lichinha",
		"Tomasinho",
		"Nachinho",
	}

	orderedList := OrderedList(items)

	expectedResult := []string{
		"1. Lichinha",
		"2. Tomasinho",
		"3. Nachinho",
	}

	for _, expected := range expectedResult {
		if !strings.Contains(orderedList, expected) {
			t.Fatalf("Missing string: %s", expected)
		}
	}
}

func TestUnorderedList(t *testing.T) {
	items := []string{
		"Lichinha",
		"Tomasinho",
		"Nachinho",
	}

	orderedList := UnorderedList(items)

	expectedResult := []string{
		"• Lichinha",
		"• Tomasinho",
		"• Nachinho",
	}

	for _, expected := range expectedResult {
		if !strings.Contains(orderedList, expected) {
			t.Fatalf("Missing string: %s", expected)
		}
	}
}

func TestEllipseText(t *testing.T) {
	testCases := []struct {
		Name                  string
		Text                  string
		MaxAmountOfCharacters int
		ExpectedText          string
	}{
		{
			Name:                  "Text with length less than max amount of characters",
			Text:                  "hola",
			MaxAmountOfCharacters: 20,
			ExpectedText:          "hola",
		},
		{
			Name:                  "Ellipse text correctly",
			Text:                  "hola que tal tu como estas? dime si eres feliz",
			MaxAmountOfCharacters: 31,
			ExpectedText:          "hola que tal tu como estas? dim...",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ellipsis := EllipseText(testCase.Text, testCase.MaxAmountOfCharacters)
			assert.Equal(t, testCase.ExpectedText, ellipsis)
		})
	}
}
