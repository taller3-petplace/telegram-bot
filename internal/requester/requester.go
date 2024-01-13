package requester

import (
	"encoding/json"
	"net/http"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/utils"
)

const (
	getPets        = "get_pets"
	registerPet    = "register_pet"
	configFilePath = "internal/config/config.json"
)

type Requester struct {
	PetsService       config.ServiceEndpoints `json:"pets_service"`
	TreatmentsService config.ServiceEndpoints `json:"treatments_service"`
	UsersService      config.ServiceEndpoints `json:"users_service"`
	clientHttp        http.Client
}

func NewRequester() (*Requester, error) {
	rawFileData, err := utils.ReadFileWithPath(configFilePath, "requester.go")
	if err != nil {
		return nil, err
	}

	var requester Requester
	err = json.Unmarshal(rawFileData, &requester)
	if err != nil {
		return nil, err
	}

	return &requester, nil
}
