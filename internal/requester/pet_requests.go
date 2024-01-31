package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils/urlutils"
)

const (
	getPets          = "get_pets"
	registerPet      = "register_pet"
	getPetByID       = "get_pet_by_id"
	headerTelegramID = "X-Telegram-Id"
)

func (r *Requester) GetPetsByOwnerID(ownerID int64) ([]domain.PetData, error) {
	operation := "GetPetsByOwnerID"
	endpointData, err := r.PetsService.GetEndpoint(getPets)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"ownerID": fmt.Sprintf("%v", ownerID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		err = fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
		logrus.Errorf("%v", err)
		return nil, err
	}

	if endpointData.QueryParams != nil {
		urlutils.AddQueryParams(request, endpointData.QueryParams.ToMap())
	}

	setTelegramHeader(request)
	request.Header.Add(headerTelegramID, fmt.Sprint(ownerID))
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing getPetsByOwnerID: %v", err)
		return nil, NewRequestError(
			fmt.Errorf("%w %s", errPerformingRequest, operation),
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response == nil {
		logrus.Error("nil response from pets service")
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return nil, errorResponse
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("error from pets service: %v", err)
		return nil, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading pets body: %v", err)
		return nil, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var petsResponse domain.PetsResponse
	err = json.Unmarshal(responseBody, &petsResponse)
	if err != nil {
		logrus.Errorf("error unmarshalling pets data: %v", err)
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingMultiplePetsData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return petsResponse.PetsData, nil
}

// RegisterPet request to register the pet of a given user
func (r *Requester) RegisterPet(petDataRequest domain.PetRequest) error {
	operation := "RegisterPet"
	endpointData, err := r.PetsService.GetEndpoint(registerPet)
	if err != nil {
		logrus.Errorf("%v", err)
		return fmt.Errorf("%w: %s", errEndpointDoesNotExist, registerPet)
	}

	url := endpointData.GetURL()
	rawBody, err := json.Marshal(petDataRequest)
	if err != nil {
		logrus.Errorf("error marshalling pet request: %v", err)
		return fmt.Errorf("%w: %v", errMarshallingPetRequest, err)
	}

	requestBody := bytes.NewReader(rawBody)
	request, err := http.NewRequest(endpointData.Method, url, requestBody)
	if err != nil {
		logrus.Errorf("error creating registerPet request: %v", err)
		return fmt.Errorf("%w: %v", errCreatingRequest, err)
	}

	setTelegramHeader(request)
	request.Header.Add(headerTelegramID, fmt.Sprint(petDataRequest.OwnerID))
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing registerPet request: %v", err)
		return NewRequestError(
			fmt.Errorf("%w %s", errPerformingRequest, operation),
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response == nil {
		logrus.Errorf("%v in registerPet", errNilResponse)
		return NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("%v", err)
		return NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	return nil
}

// GetPetData fetch information about a pet based on the given ID
func (r *Requester) GetPetData(petID int) (domain.PetData, error) {
	operation := "GetPetData"
	endpointData, err := r.PetsService.GetEndpoint(getPetByID)
	if err != nil {
		logrus.Errorf("%v", err)
		return domain.PetData{}, fmt.Errorf("%w: %s", errEndpointDoesNotExist, getPetByID)
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		logrus.Errorf("error creating getPet request: %v", err)
		return domain.PetData{}, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing getPet request: %v", err)
		return domain.PetData{}, NewRequestError(
			fmt.Errorf("%w %s", errPerformingRequest, operation),
			http.StatusInternalServerError,
			err.Error(),
		)
	}

	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response == nil {
		logrus.Errorf("%v getting pet data", errNilResponse)
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return domain.PetData{}, errorResponse
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("%v", err)
		return domain.PetData{}, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading petData body: %v", err)
		return domain.PetData{}, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var petData domain.PetData
	err = json.Unmarshal(responseBody, &petData)
	if err != nil {
		logrus.Errorf("error unmarshalling pet data: %v", err)
		return domain.PetData{}, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingPetData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return petData, nil
}
