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

	for _, notificationToSend := range notifications {
		// Best effort
		telegramID, err := strconv.Atoi(notificationToSend.TelegramID)
		if err != nil {
			logrus.Errorf("error invalid telegramID: %s", notificationToSend.TelegramID)
			continue
		}

		err = ns.telegramBot.SendNotification(int64(telegramID), notificationToSend.Message)
		if err != nil {
			logrus.Errorf("error sending notification, telegram_id: %s", notificationToSend.TelegramID)
		}
	}

	c.JSON(http.StatusOK, nil)
}
