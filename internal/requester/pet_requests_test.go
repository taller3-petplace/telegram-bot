package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/requester/internal/mock"
	"testing"
	"time"
)

const (
	ownerID     = int64(69)
	testBaseURL = "https://test"
)

type expectedServiceConfig struct {
	BaseURL           string
	ExpectedEndpoints map[string]config.Endpoint
}

func TestNewRequester(t *testing.T) {
	client := http.Client{}
	requester, err := NewRequester(&client)
	assert.NoError(t, err)

	expectedPetsServiceConfig := expectedServiceConfig{
		BaseURL:           "http://localhost:8712/pets",
		ExpectedEndpoints: getExpectedPetsServiceEndpoints(),
	}
	assertServiceConfig(t, requester.PetsService, expectedPetsServiceConfig)

	expectedTreatmentsServiceConfig := expectedServiceConfig{
		BaseURL:           "http://localhost:8712/treatments",
		ExpectedEndpoints: getExpectedTreatmentsServiceEndpoints(),
	}
	assertServiceConfig(t, requester.TreatmentsService, expectedTreatmentsServiceConfig)

	expectedUsersServiceConfig := expectedServiceConfig{
		BaseURL:           "http://localhost:8712/users",
		ExpectedEndpoints: getExpectedUsersServiceEndpoints(),
	}
	assertServiceConfig(t, requester.UsersService, expectedUsersServiceConfig)
}

func assertServiceConfig(t *testing.T, service config.ServiceEndpoints, expectedResults expectedServiceConfig) {
	assert.Equal(t, expectedResults.BaseURL, service.Base)
	assert.Equal(t, len(expectedResults.ExpectedEndpoints), len(service.Endpoints))

	for endpointAlias, endpointData := range service.Endpoints {
		expectedEndpoint, found := expectedResults.ExpectedEndpoints[endpointAlias]
		if !found {
			t.Fatalf("endpoint %s is missing", endpointAlias)
		}

		assert.Equal(t, expectedEndpoint.Path, endpointData.Path)
		assert.Equal(t, expectedEndpoint.Method, endpointData.Method)
		assert.Equal(t, expectedEndpoint.QueryParams, endpointData.QueryParams)
		expectedURL := service.Base + expectedEndpoint.Path
		assert.Equal(t, expectedURL, endpointData.GetURL())
	}
}

func getExpectedPetsServiceEndpoints() map[string]config.Endpoint {
	return map[string]config.Endpoint{
		"register_pet": {
			Path:   "/pet",
			Method: http.MethodPost,
		},
		"get_pets": {
			Path:   "/owner/{ownerID}",
			Method: http.MethodGet,
			QueryParams: &config.QueryParams{
				Offset: 0,
				Limit:  100,
			},
		},
		"get_pet_by_id": {
			Path:   "/pet/{petID}",
			Method: http.MethodGet,
		},
	}
}

func getExpectedTreatmentsServiceEndpoints() map[string]config.Endpoint {
	return map[string]config.Endpoint{
		"get_pet_treatments": {
			Path:   "/treatment/pet/{petID}",
			Method: http.MethodGet,
			QueryParams: &config.QueryParams{
				Offset: 0,
				Limit:  5,
			},
		},
		"get_treatment": {
			Path:   "/treatment/{treatmentID}",
			Method: http.MethodGet,
		},
	}
}

func getExpectedUsersServiceEndpoints() map[string]config.Endpoint {
	return map[string]config.Endpoint{
		"user_fetcher": {
			Path:   "/telegram_id/{telegramID}",
			Method: http.MethodGet,
		},
	}
}

type clientMockConfig struct {
	RequestBody  io.Reader
	ResponseBody *http.Response
	Err          error
}

func TestRequesterGetPetsByOwnerID(t *testing.T) {
	petsServiceEndpoints := getExpectedPetsServiceEndpoints()
	getPetsByOwnerIDEndpoint := petsServiceEndpoints[getPets]
	getPetsByOwnerIDEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := petsServiceEndpoints[getPets]
	invalidEndpoint.Method = "hola que tal tu como estas? dime si eres feliz"

	requester := Requester{
		PetsService: config.ServiceEndpoints{
			Endpoints: petsServiceEndpoints,
		},
	}

	petsServiceError := petServiceErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "error cae el soooool en tu balcoooooon",
	}
	serviceErrorRaw, err := json.Marshal(petsServiceError)
	require.NoError(t, err)

	petsData := []domain.PetDataIdentifier{
		{
			ID:   1,
			Name: "Cartucho",
			Type: "DOG",
		},
		{
			ID:   2,
			Name: "Pantufla",
			Type: "CAT",
		},
	}

	rawPetsData, err := json.Marshal(petsData)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		Requester        Requester
		ClientMockConfig *clientMockConfig
		ExpectsError     bool
		ExpectedError    error
		ExpectedPetsData []domain.PetDataIdentifier
	}{
		{
			Name: "Endpoint does not exist",
			Requester: Requester{
				PetsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{}},
			},
			ExpectsError:  true,
			ExpectedError: errEndpointDoesNotExist,
		},
		{
			Name: "Error creating request",
			Requester: Requester{
				PetsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{
					getPets: invalidEndpoint,
				}},
			},
			ExpectsError:  true,
			ExpectedError: errCreatingRequest,
		},
		{
			Name:      "Error performing request",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: nil,
				Err:          fmt.Errorf("internal error performing request"),
			},
			ExpectsError:  true,
			ExpectedError: errPerformingRequest,
		},
		{
			Name:      "Error nil response",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: nil,
				Err:          nil,
			},
			ExpectsError:  true,
			ExpectedError: errNilResponse,
		},
		{
			Name:      "Error from pets service",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(serviceErrorRaw)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: fmt.Errorf(petsServiceError.GetMessage()),
		},
		{
			Name:      "Error unmarshalling pets data",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id": "69abc"}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingPetsData,
		},
		{
			Name:      "Get pets data correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawPetsData)),
				},
				Err: nil,
			},
			ExpectsError:     false,
			ExpectedPetsData: petsData,
			ExpectedError:    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			clientMock := mock.NewMockhttpClienter(gomock.NewController(t))
			if testCase.ClientMockConfig != nil {
				clientMock.EXPECT().
					Do(gomock.Any()).
					Return(testCase.ClientMockConfig.ResponseBody, testCase.ClientMockConfig.Err)
			}

			testCase.Requester.clientHTTP = clientMock

			petsDataResponse, err := testCase.Requester.GetPetsByOwnerID(ownerID)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedPetsData, petsDataResponse)
		})
	}
}

func TestRequesterRegisterPet(t *testing.T) {
	petsServiceEndpoints := getExpectedPetsServiceEndpoints()
	registerPetEndpoint := petsServiceEndpoints[registerPet]
	registerPetEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := petsServiceEndpoints[registerPet]
	invalidEndpoint.Method = "hola que tal tu como estas? dime si eres feliz"

	requester := Requester{
		PetsService: config.ServiceEndpoints{
			Endpoints: petsServiceEndpoints,
		},
	}

	petsServiceError := petServiceErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "error alla le estan registrando una mascota",
	}
	serviceErrorRaw, err := json.Marshal(petsServiceError)
	require.NoError(t, err)

	petRequest := domain.PetRequest{
		Name:         "Turron",
		Type:         "DOG",
		RegisterDate: time.Now(),
		BirthDate:    "2013/05/15",
		OwnerID:      ownerID,
	}
	rawPetRequest, err := json.Marshal(petRequest)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		Requester        Requester
		ClientMockConfig *clientMockConfig
		ExpectsError     bool
		ExpectedError    error
	}{
		{
			Name: "Endpoint does not exist",
			Requester: Requester{
				PetsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{}},
			},
			ExpectsError:  true,
			ExpectedError: errEndpointDoesNotExist,
		},
		{
			Name: "Error creating request",
			Requester: Requester{
				PetsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{
					registerPet: invalidEndpoint,
				}},
			},
			ExpectsError:  true,
			ExpectedError: errCreatingRequest,
		},
		{
			Name:      "Error performing request",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: nil,
				Err:          fmt.Errorf("internal error performing request"),
			},
			ExpectsError:  true,
			ExpectedError: errPerformingRequest,
		},
		{
			Name:      "Error nil response",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: nil,
				Err:          nil,
			},
			ExpectsError:  true,
			ExpectedError: errNilResponse,
		},
		{
			Name:      "Error from pets service",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(serviceErrorRaw)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: fmt.Errorf(petsServiceError.GetMessage()),
		},
		{
			Name:      "Register pet correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				RequestBody: bytes.NewReader(rawPetRequest),
				ResponseBody: &http.Response{
					StatusCode: http.StatusCreated,
					Body:       nil,
				},
				Err: nil,
			},
			ExpectsError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			clientMock := mock.NewMockhttpClienter(gomock.NewController(t))
			if testCase.ClientMockConfig != nil {
				clientMock.EXPECT().
					Do(gomock.Any()).
					Return(testCase.ClientMockConfig.ResponseBody, testCase.ClientMockConfig.Err)
			}

			testCase.Requester.clientHTTP = clientMock

			err := testCase.Requester.RegisterPet(petRequest)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
