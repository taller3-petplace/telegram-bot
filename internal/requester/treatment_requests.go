package requester

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils"
	"telegram-bot/internal/utils/urlutils"
)

const (
	getPetTreatments = "get_pet_treatments"
	getTreatment     = "get_treatment"
	getVaccines      = "get_vaccines"
)

// GetTreatmentsByPetID fetches all treatments of the given pet.
// The treatments are ordered from most recent to oldest
func (r *Requester) GetTreatmentsByPetID(petID int) ([]domain.Treatment, error) {
	operation := "GetTreatmentsByPetID"
	endpointData, err := r.TreatmentsService.GetEndpoint(getPetTreatments)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		logrus.Errorf("error creating getTreatmentsByPetID request: %v", err)
		return nil, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	if endpointData.QueryParams != nil {
		urlutils.AddQueryParams(request, endpointData.QueryParams.ToMap())
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing getTreatmentsByPetID: %v", err)
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
		logrus.Errorf("%v", errNilResponse)
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return nil, errorResponse
	}

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading treatments body: %v", err)
		return nil, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var petTreatments []domain.Treatment
	err = json.Unmarshal(responseBody, &petTreatments)
	if err != nil {
		logrus.Errorf("error unmarshallin pet treatments: %v", err)
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingMultipleTreatments, err),
			http.StatusInternalServerError,
			"",
		)
	}

	utils.SortElementsByDate(petTreatments)
	return petTreatments, nil
}

// GetTreatment fetches all the information about the given treatment
func (r *Requester) GetTreatment(treatmentID string) (domain.Treatment, error) {
	operation := "GetTreatment"
	endpointData, err := r.TreatmentsService.GetEndpoint(getTreatment)
	if err != nil {
		logrus.Errorf("%v", err)
		return domain.Treatment{}, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"treatmentID": fmt.Sprintf("%v", treatmentID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		logrus.Errorf("error creting getTreatment request: %v", err)
		return domain.Treatment{}, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing getTreatment request: %v", err)
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
		logrus.Errorf("%v", errNilResponse)
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return domain.Treatment{}, errorResponse
	}

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("error from treatments service: %v", err)
		return domain.Treatment{}, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading treatment body: %v", err)
		return domain.Treatment{}, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var treatmentData domain.Treatment
	err = json.Unmarshal(responseBody, &treatmentData)
	if err != nil {
		logrus.Errorf("error unmarshalling treatment data: %v", err)
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
		logrus.Errorf("%v", err)
		return nil, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"petID": fmt.Sprintf("%v", petID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		logrus.Errorf("error creating getVaccines request: %v", err)
		return nil, fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing getVaccines request: %v", err)
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
		logrus.Errorf("%v in getVaccines", errNilResponse)
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return nil, errorResponse
	}

	err = ErrPolicyFunc[treatmentServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("error from treatments service %v", err)
		return nil, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading vaccines body: %v", err)
		return nil, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var vaccinesResponse []domain.VaccineResponse
	err = json.Unmarshal(responseBody, &vaccinesResponse)
	if err != nil {
		logrus.Errorf("error unmarshalling vaccines response: %v", err)
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingVaccinesData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	// Group vaccines by name
	vaccinesMap := make(map[string][]domain.VaccineResponse)
	for _, vac := range vaccinesResponse {
		vaccinesMap[vac.Name] = append(vaccinesMap[vac.Name], vac)
	}

	// Generate response
	var vaccines []domain.Vaccine
	for vaccineName, vaccinesApplied := range vaccinesMap {
		utils.SortElementsByDate(vaccinesApplied)

		vaccine := domain.Vaccine{
			Name:          vaccineName,
			AmountOfDoses: len(vaccinesApplied),
			FirstDose:     vaccinesApplied[len(vaccinesApplied)-1].Date,
			LastDose:      vaccinesApplied[0].Date,
		}
		vaccines = append(vaccines, vaccine)
	}

	utils.SortElementsByDate(vaccines)

	return vaccines, nil
}
