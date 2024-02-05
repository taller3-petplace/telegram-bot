package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
	"telegram-bot/internal/bot/internal/template"
	"time"
)

const editTTL = 3 * time.Minute

// help sends the commands that botName supports
func (tb *TelegramBot) help(c tele.Context) error {
	message := fmt.Sprintf("Available commands:\n\n%s", template.Commands())
	return c.Send(message)
}

// editMessageHandler handle message edition. If the message is edited after than editTTL, the bot does not perform any action
func (tb *TelegramBot) editMessageHandler(c tele.Context) error {
	currentTime := c.Message().Time()

	if !c.Message().LastEdited().Before(currentTime.Add(editTTL)) {
		return c.Send("The time to edit the message has expired, start again")
	}

	return tb.textHandler(c)
}

// textHandler handles text input from the user. If the text contains a specific endpoint, it forwards the message to the corresponding handler
func (tb *TelegramBot) textHandler(c tele.Context) error {
	message := c.Message().Text
	if strings.Contains(message, registerPetEndpoint) {
		return tb.createPetRecord(c)
	}

	if strings.Contains(message, registerNotificationEndpoint) {
		return tb.registerNotification(c)
	}

	return c.Send("I don't understand your input, execute /help to check what can I do for you")
}
