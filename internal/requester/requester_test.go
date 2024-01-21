package requester

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"telegram-bot/internal/requester/internal/config"
	"testing"
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
		BaseURL:           "https://api.lnt.digital/pets",
		ExpectedEndpoints: getExpectedPetsServiceEndpoints(),
	}
	assertServiceConfig(t, requester.PetsService, expectedPetsServiceConfig)

	expectedTreatmentsServiceConfig := expectedServiceConfig{
		BaseURL:           "https://api.lnt.digital/treatments",
		ExpectedEndpoints: getExpectedTreatmentsServiceEndpoints(),
	}
	assertServiceConfig(t, requester.TreatmentsService, expectedTreatmentsServiceConfig)

	expectedUsersServiceConfig := expectedServiceConfig{
		BaseURL:           "https://api.lnt.digital/users",
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
			Path:   "/treatment/specific/{treatmentID}",
			Method: http.MethodGet,
		},
		"get_vaccines": {
			Path:   "/application/pet/{petID}",
			Method: http.MethodGet,
			QueryParams: &config.QueryParams{
				Offset: 0,
				Limit:  100,
			},
		},
	}
}

func getExpectedUsersServiceEndpoints() map[string]config.Endpoint {
	return map[string]config.Endpoint{
		"get_user": {
			Path:   "/telegram_id/{telegramID}",
			Method: http.MethodGet,
		},
	}
}
