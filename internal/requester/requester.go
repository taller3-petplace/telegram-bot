package requester

import (
	"gopkg.in/yaml.v3"
	"telegram-bot/internal/utils"
)

type Requester struct {
	PetsService       serviceEndpoints `yaml:"pets_service"`
	TreatmentsService serviceEndpoints `yaml:"treatments_service"`
	UsersService      serviceEndpoints `yaml:"users_service"`
}

type serviceEndpoints struct {
	Base      string              `yaml:"base"`
	Endpoints map[string]endpoint `yaml:"endpoints"`
}

type endpoint struct {
	Path        string       `yaml:"path"`
	QueryParams *queryParams `yaml:"query_params"`
}

type queryParams struct {
	Offset int `yaml:"offset"`
	Limit  int `yaml:"limit"`
}

func NewRequester() (*Requester, error) {
	rawFileData, err := utils.ReadFileWithPath("internal/config/config.yml", "requester.go")
	if err != nil {
		return nil, err
	}

	var requester Requester
	err = yaml.Unmarshal(rawFileData, &requester)
	if err != nil {
		return nil, err
	}

	return &requester, nil
}
