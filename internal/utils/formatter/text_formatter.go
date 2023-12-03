package formatter

import "fmt"

func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func Italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

func Link(text string, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
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
