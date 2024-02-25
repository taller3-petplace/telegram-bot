package requester

import (
	"encoding/json"
	"net/http"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/utils"
)

const (
	configFilePath = "internal/requester/internal/config/config.json"
	telegramHeader = "X-Telegram-App"
)

type httpClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

type Requester struct {
	PetsService          config.ServiceEndpoints `json:"pets_service"`
	TreatmentsService    config.ServiceEndpoints `json:"treatments_service"`
	UsersService         config.ServiceEndpoints `json:"users_service"`
	NotificationsService config.ServiceEndpoints `json:"notifications_service"`
	clientHTTP           httpClienter
}

func NewRequester(client httpClienter) (*Requester, error) {
	rawFileData, err := utils.ReadFileWithPath(configFilePath, "requester.go")
	if err != nil {
		return nil, err
	}

	var requester Requester
	err = json.Unmarshal(rawFileData, &requester)
	if err != nil {
		return nil, err
	}

	requester.clientHTTP = client

	return &requester, nil
}

// setTelegramHeader sets a header to indicate that the request come from Telegram Service
func setTelegramHeader(request *http.Request) {
	request.Header.Add(telegramHeader, "true")
}
