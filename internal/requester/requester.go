package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/urlutils"
	"time"
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
	err = yaml.Unmarshal(rawFileData, &requester)
	if err != nil {
		return nil, err
	}

	return &requester, nil
}

func (r *Requester) GetPetsByOwnerID(ownerID int64) ([]domain.PetDataIdentifier, error) {
	operation := "GetPetsByOwnerID"
	endpointData, endpointExists := r.PetsService.Endpoints[getPets]
	if !endpointExists {
		return nil, fmt.Errorf("%w: %s", errEndpointDoesNotExist, getPets)
	}

	// ToDo: perform this in Unmarshall
	url := r.PetsService.Base + endpointData.Path
	url = urlutils.FormatURL(url, map[string]string{"owner_id": fmt.Sprintf("%v", ownerID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error doing getPets request: %v", err)
	}

	urlutils.AddQueryParams(request, endpointData.QueryParams.ToMap())

	r.clientHttp.Timeout = 5 * time.Second

	response, err := r.clientHttp.Do(request)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w %s", errPerformingRequest, operation),
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response == nil {
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return nil, errorResponse
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		return nil, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var petsData []domain.PetDataIdentifier
	err = json.Unmarshal(responseBody, &petsData)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingPetsData, err),
			response.StatusCode,
			"",
		)
	}

	return petsData, nil
}

func (r *Requester) RegisterPet(petDataRequest domain.PetRequest) error {
	endpointData, endpointExists := r.PetsService.Endpoints[registerPet]
	if !endpointExists {
		return fmt.Errorf("%w: %s", errEndpointDoesNotExist, registerPet)
	}

	url := r.PetsService.Base + endpointData.Path
	rawBody, err := json.Marshal(petDataRequest)
	if err != nil {
		return fmt.Errorf("%w: %v", errMarshallingPetRequest, err)
	}

	requestBody := bytes.NewReader(rawBody)
	request, err := http.NewRequest(endpointData.Method, url, requestBody)
	if err != nil {
		return fmt.Errorf("%w: %v", errCreatingRequest, err)
	}

	r.clientHttp.Timeout = 5 * time.Second
	response, err := r.clientHttp.Do(request)
	if err != nil {
		return NewRequestError(
			fmt.Errorf("%w RegisterPet", errPerformingRequest),
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response == nil {
		return NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			"RegisterPet",
		)
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		return NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	return nil
}
