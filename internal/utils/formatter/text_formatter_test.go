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
