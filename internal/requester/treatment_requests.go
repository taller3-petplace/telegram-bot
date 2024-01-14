package requester

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils/urlutils"
)

const (
	getPetTreatments      = "get_pet_treatments"
	getTreatment          = "get_treatment"
	getVaccines           = "get_vaccines"
	maxAmountOfTreatments = 5
)

// GetTreatmentsByPetID fetches the last 5 treatments of the given pet.
// The treatments go from new ones to old ones
func (r *Requester) GetTreatmentsByPetID(petID int) ([]domain.Treatment, error) {
	operation := "GetTreatmentsByPetID"
	endpointData, err := r.TreatmentsService.GetEndpoint(getPetTreatments)
	if err != nil {
		return nil, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
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

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
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

	var petTreatments []domain.Treatment
	err = json.Unmarshal(responseBody, &petTreatments)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingMultipleTreatments, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return petTreatments, nil
}

// GetTreatment fetches all the information about the given treatment
func (r *Requester) GetTreatment(treatmentID int) (domain.Treatment, error) {
	operation := "GetTreatment"
	endpointData, err := r.TreatmentsService.GetEndpoint(getTreatment)
	if err != nil {
		return domain.Treatment{}, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"treatmentID": fmt.Sprintf("%v", treatmentID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return domain.Treatment{}, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	response, err := r.clientHTTP.Do(request)
	if err != nil {
		return domain.Treatment{}, NewRequestError(
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
		return domain.Treatment{}, errorResponse
	}

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
	if err != nil {
		return domain.Treatment{}, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Treatment{}, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var treatmentData domain.Treatment
	err = json.Unmarshal(responseBody, &treatmentData)
	if err != nil {
		return domain.Treatment{}, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingTreatmentData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return treatmentData, nil
}

// GetVaccines fetches all the vaccines that were applied to the pet
func (r *Requester) GetVaccines(petID int) ([]domain.Vaccine, error) {
	operation := "GetVaccines"
	endpointData, err := r.TreatmentsService.GetEndpoint(getVaccines)
	if err != nil {
		return nil, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
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

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
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

	var vaccines []domain.Vaccine
	err = json.Unmarshal(responseBody, &vaccines)
	if err != nil {
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingVaccinesData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return vaccines, nil
}
