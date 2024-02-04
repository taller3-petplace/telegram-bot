package notification

type Notification struct {
	TelegramID string `json:"telegram_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}
