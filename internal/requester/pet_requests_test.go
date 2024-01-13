package requester

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/requester/internal/config"
	"telegram-bot/internal/requester/internal/mock"
	"telegram-bot/internal/utils/urlutils"
	"testing"
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
	ResponseBody *http.Response
	Err          error
}

func TestRequester_GetPetsByOwnerID(t *testing.T) {
	petsServiceEndpoints := getExpectedPetsServiceEndpoints()
	getPetsByOwnerIDEndpoint := petsServiceEndpoints[getPets]
	getPetsByOwnerIDEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := petsServiceEndpoints[getPets]
	invalidEndpoint.Method = "hola que tal tu como estas? dime si eres feliz"

	expectedURL := getPetsByOwnerIDEndpoint.GetURL()
	expectedURL = urlutils.FormatURL(expectedURL, map[string]string{"ownerID": fmt.Sprint(ownerID)})

	requester := Requester{
		PetsService: config.ServiceEndpoints{
			Endpoints: petsServiceEndpoints,
		},
	}

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

			petsData, err := testCase.Requester.GetPetsByOwnerID(ownerID)
			if testCase.ExpectsError {
				assert.ErrorIs(t, err, testCase.ExpectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedPetsData, petsData)
		})
	}
}
