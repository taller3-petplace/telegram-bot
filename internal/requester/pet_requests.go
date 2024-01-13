package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils/urlutils"
)

const (
	getPets     = "get_pets"
	registerPet = "register_pet"
	getPetByID  = "get_pet_by_id"
)

func (r *Requester) GetPetsByOwnerID(ownerID int64) ([]domain.PetData, error) {
	operation := "GetPetsByOwnerID"
	endpointData, endpointExists := r.PetsService.Endpoints[getPets]
	if !endpointExists {
		return nil, fmt.Errorf("%w: %s", errEndpointDoesNotExist, getPets)
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"ownerID": fmt.Sprintf("%v", ownerID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	if endpointData.QueryParams != nil {
		urlutils.AddQueryParams(request, endpointData.QueryParams.ToMap())
	}

	response, err := r.clientHTTP.Do(request)
	if err != nil {
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

	var petsData []domain.PetData
	err = json.Unmarshal(responseBody, &petsData)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingMultiplePetsData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return petsData, nil
}

// RegisterPet request to register the pet of a given user
func (r *Requester) RegisterPet(petDataRequest domain.PetRequest) error {
	operation := "RegisterPet"
	endpointData, endpointExists := r.PetsService.Endpoints[registerPet]
	if !endpointExists {
		return fmt.Errorf("%w: %s", errEndpointDoesNotExist, registerPet)
	}

	url := endpointData.GetURL()
	rawBody, err := json.Marshal(petDataRequest)
	if err != nil {
		return fmt.Errorf("%w: %v", errMarshallingPetRequest, err)
	}

	requestBody := bytes.NewReader(rawBody)
	request, err := http.NewRequest(endpointData.Method, url, requestBody)
	if err != nil {
		return fmt.Errorf("%w: %v", errCreatingRequest, err)
	}

	response, err := r.clientHTTP.Do(request)
	if err != nil {
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
		return NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
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

// GetPetData fetch information about a pet based on the given ID
func (r *Requester) GetPetData(petID int) (domain.PetData, error) {
	operation := "GetPetData"
	endpointData, endpointExists := r.PetsService.Endpoints[getPetByID]
	if !endpointExists {
		return domain.PetData{}, fmt.Errorf("%w: %s", errEndpointDoesNotExist, getPetByID)
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return domain.PetData{}, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	response, err := r.clientHTTP.Do(request)
	if err != nil {
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
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return domain.PetData{}, errorResponse
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		return domain.PetData{}, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.PetData{}, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var petData domain.PetData
	err = json.Unmarshal(responseBody, &petData)
	if err != nil {
		return domain.PetData{}, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingPetData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return petData, nil
}
