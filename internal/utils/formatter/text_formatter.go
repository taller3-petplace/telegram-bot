package formatter

import (
	"fmt"
	"strings"
)

func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func Italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

func Link(text string, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

func Capitalize(input string) string {
	firstLetter := strings.ToUpper(input[0:1])
	return firstLetter + input[1:]
}

func OrderedList(items []string) string {
	var output string
	for idx, item := range items {
		output += fmt.Sprintf("%d. %s\n\n", idx+1, item)
	}

	return output
}

func UnorderedList(items []string) string {
	var output string
	for _, item := range items {
		output += fmt.Sprintf("\t• %s\n\n", item)
	}

	return output
}

func EllipseText(text string, maxAmountOfCharacters int) string {
	if len(text) < maxAmountOfCharacters {
		return text
	}

	text = text[:maxAmountOfCharacters]
	return text + "..."
}

func UnderlineText(text string) string {
	return fmt.Sprintf("__%s__", text)
}

func SpoilerText(text string) string {
	return fmt.Sprintf("||%s||", text)
}
