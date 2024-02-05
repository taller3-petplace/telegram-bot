package sender

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/sender/internal/notification"
)

type NotificationsSender struct {
	telegramBot *bot.TelegramBot
}

func NewNotificationSender(telegramBot *bot.TelegramBot) *NotificationsSender {
	return &NotificationsSender{
		telegramBot: telegramBot,
	}
}

type summary struct {
	OK   int `json:"ok"`
	Fail int `json:"fail"`
}

// TriggerNotifications sends each notification that receives to the corresponding user. Best effort procedure
func (ns *NotificationsSender) TriggerNotifications(c *gin.Context) {
	var notifications []notification.Notification
	err := c.ShouldBindJSON(&notifications)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("error unmarshaling body: %v", err.Error()),
		})
		return
	}

	counter := 0
	for _, notificationToSend := range notifications {
		// Best effort
		telegramID, err := strconv.Atoi(notificationToSend.TelegramID)
		if err != nil {
			logrus.Errorf("error invalid telegramID %s: %v", notificationToSend.TelegramID, err)
			continue
		}

		err = ns.telegramBot.SendNotification(int64(telegramID), notificationToSend.Message)
		if err != nil {
			logrus.Errorf("error sending notification, telegram_id: %s: %v", notificationToSend.TelegramID, err)
			continue
		}
		counter++
	}

	c.JSON(http.StatusOK, summary{
		OK:   counter,
		Fail: len(notifications) - counter,
	})
}
