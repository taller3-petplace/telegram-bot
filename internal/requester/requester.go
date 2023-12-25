package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"regexp"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils"
	"time"
)

const (
	getPets     = "get_pets"
	registerPet = "register_pet"
)

type Requester struct {
	PetsService       serviceEndpoints `yaml:"pets_service"`
	TreatmentsService serviceEndpoints `yaml:"treatments_service"`
	UsersService      serviceEndpoints `yaml:"users_service"`
	clientHttp        http.Client
}

type serviceEndpoints struct {
	Base      string              `yaml:"base"`
	Endpoints map[string]endpoint `yaml:"endpoints"`
}

type endpoint struct {
	Path        string       `yaml:"path"`
	Method      string       `yaml:"method"`
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

func (r *Requester) GetPetsByOwnerID(ownerID int64) ([]domain.PetDataSummary, error) {
	operation := "GetPetsByOwnerID"
	endpointData, endpointExists := r.PetsService.Endpoints[getPets]
	if !endpointExists {
		return nil, fmt.Errorf("%w: %s", errEndpointDoesNotExist, getPets)
	}

	url := r.PetsService.Base + endpointData.Path
	url = FormatURL(url, map[string]string{"owner_id": fmt.Sprintf("%v", ownerID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error doing getPets request: %v", err)
	}

	// ToDo: function to add query params
	queryParamsValues := request.URL.Query()
	queryParamsValues.Add("limit", fmt.Sprintf("%v", endpointData.QueryParams.Limit))
	queryParamsValues.Add("offset", fmt.Sprintf("%v", endpointData.QueryParams.Offset))
	request.URL.RawQuery = queryParamsValues.Encode()

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

	var petsData domain.PetsData
	err = json.Unmarshal(responseBody, &petsData)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingPetsData, err),
			response.StatusCode,
			"",
		)
	}

	return petsData.PetsData, nil
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

// FormatURL formats the given URL setting the values from the map to it. Is not-in-place
func FormatURL(url string, params map[string]string) string {
	formattedURL := url
	for param, value := range params {
		regex := regexp.MustCompile(fmt.Sprintf("{%s}", param))
		formattedURL = regex.ReplaceAllString(formattedURL, value)
	}

	return formattedURL
}
