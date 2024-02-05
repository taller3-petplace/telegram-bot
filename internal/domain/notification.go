package domain

import "time"

const via = "telegram"

type NotificationRequest struct {
	TelegramID string     `json:"telegram_id"`
	Via        string     `json:"via"`
	Message    string     `json:"message"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	Hours      []string   `json:"hours"`
}

func NewNotificationRequest(telegramID string, sms string, startDate time.Time, endDate *time.Time, hours []string) NotificationRequest {
	return NotificationRequest{
		TelegramID: telegramID,
		Via:        via,
		Message:    sms,
		StartDate:  startDate,
		EndDate:    endDate,
		Hours:      hours,
	}

}

type NotificationResponse struct {
	ID        string     `json:"id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Hour      string     `json:"hour"`
}
