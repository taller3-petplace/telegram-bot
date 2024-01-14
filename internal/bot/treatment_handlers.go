package bot

import (
	tele "gopkg.in/telebot.v3"
)

func (tb *TelegramBot) showVaccines(c tele.Context) error {
	return c.Send("implement me")
}

func (tb *TelegramBot) medicalHistory(c tele.Context) error {
	return c.Send("implement me")
}
