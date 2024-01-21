package requester

import (
	"encoding/json"
	"net/http"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/utils"
)

const (
	configFilePath = "internal/requester/internal/config/config.json"
)

type httpClienter interface {
	Do(req *http.Request) (*http.Response, error)
}

type Requester struct {
	PetsService       config.ServiceEndpoints `json:"pets_service"`
	TreatmentsService config.ServiceEndpoints `json:"treatments_service"`
	UsersService      config.ServiceEndpoints `json:"users_service"`
	clientHTTP        httpClienter
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
