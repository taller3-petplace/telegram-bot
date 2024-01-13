package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceEndpoints_GetEndpoint(t *testing.T) {
	baseURL := "https://music.youtube.com"
	existentEndpoint := "temaiken"
	expectedEndpoint := Endpoint{
		Path:   "/watch?v=b6Zc4nUgxHk&si=G0wQYhjFfYnnY2bT",
		Method: "GET",
	}

	nonExistentEndpoint := "sarasa"
	serviceEndpoints := ServiceEndpoints{
		Base: baseURL,
		Endpoints: map[string]Endpoint{
			existentEndpoint: expectedEndpoint,
		},
	}

	t.Run("Get endpoint correctly", func(t *testing.T) {
		endpointData, err := serviceEndpoints.GetEndpoint(existentEndpoint)
		assert.NoError(t, err)
		assert.Equal(t, expectedEndpoint, endpointData)
	})

	t.Run("Endpoint does not exist", func(t *testing.T) {
		_, err := serviceEndpoints.GetEndpoint(nonExistentEndpoint)
		assert.ErrorIs(t, err, errEndpointDoesNotExist)
	})
}

func TestServiceEndpoints_UnmarshalJSON(t *testing.T) {
	baseURL := "https://www.youtube.com"
	path := "/watch?v=RWIJExat-lQ&ab_channel=Estoes%C2%A1FA%21"
	endpointAlias := "hitazo"
	// Input test case
	inputServiceEndpointsData := ServiceEndpoints{
		Base: baseURL,
		Endpoints: map[string]Endpoint{
			endpointAlias: {
				Path:   path,
				Method: "GET",
				QueryParams: &QueryParams{
					Limit: defaultLimit,
				},
			},
		},
	}

	rawServiceEndpointsData, err := json.Marshal(inputServiceEndpointsData)
	require.NoError(t, err)

	// Expected results
	expectedServiceEndpoints := ServiceEndpoints{
		Base: baseURL,
		Endpoints: map[string]Endpoint{
			endpointAlias: {
				Path:   path,
				Method: "GET",
				QueryParams: &QueryParams{
					Limit:  defaultLimit,
					Offset: 0,
				},
				baseURL: baseURL,
			},
		},
	}

	testCases := []struct {
		Name                     string
		RawServiceData           []byte
		ExpectsError             bool
		ExpectedServiceEndpoints ServiceEndpoints
		ExpectedError            error
	}{
		{
			Name:           "Error due to ServiceEndpoints data is missing",
			RawServiceData: nil,
			ExpectsError:   true,
			ExpectedError:  errServiceEndpointDataMissing,
		},
		{
			Name:           "Error unmarshalling ServiceEndpoints data",
			RawServiceData: []byte(`{`),
			ExpectsError:   true,
			ExpectedError:  errUnmarshallingServiceData,
		},
		{
			Name:                     "ServiceEndpoints data unmarshalled correctly",
			RawServiceData:           rawServiceEndpointsData,
			ExpectsError:             false,
			ExpectedServiceEndpoints: expectedServiceEndpoints,
			ExpectedError:            nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			serviceEndpoints := ServiceEndpoints{}
			err := serviceEndpoints.UnmarshalJSON(testCase.RawServiceData)

			if testCase.ExpectsError {
				assert.ErrorIs(t, err, testCase.ExpectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedServiceEndpoints, serviceEndpoints)
		})
	}
}

func TestQueryParams_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		Name                string
		RawData             []byte
		ExpectsError        bool
		ExpectedQueryParams QueryParams
		ExpectedError       error
	}{
		{
			Name:                "Error unmarshalling query params data",
			RawData:             []byte("{wrong json format"),
			ExpectsError:        true,
			ExpectedQueryParams: QueryParams{},
			ExpectedError:       errUnmarshallingQueryParams,
		},
		{
			Name:                "Empty raw data",
			RawData:             nil,
			ExpectsError:        false,
			ExpectedQueryParams: QueryParams{},
			ExpectedError:       nil,
		},
		{
			Name:                "If limit is missing, a default one is set",
			RawData:             []byte(`{"offset": 5}`),
			ExpectsError:        false,
			ExpectedQueryParams: QueryParams{Limit: defaultLimit, Offset: 5},
			ExpectedError:       nil,
		},
		{
			Name:                "Query params raw data unmarshalled correctly",
			RawData:             []byte(`{"limit": 10,"offset": 5}`),
			ExpectsError:        false,
			ExpectedQueryParams: QueryParams{Limit: 10, Offset: 5},
			ExpectedError:       nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			queryParams := QueryParams{}
			err := queryParams.UnmarshalJSON(testCase.RawData)
			if testCase.ExpectsError {
				assert.ErrorIs(t, err, testCase.ExpectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedQueryParams, queryParams)
		})
	}
}

func TestQueryParams_ToMap(t *testing.T) {
	queryParams := QueryParams{
		Limit:  defaultLimit,
		Offset: 5,
	}

	queryParamsMap := queryParams.ToMap()
	assert.Equal(t, queryParamsMap["limit"], fmt.Sprint(defaultLimit))
	assert.Equal(t, queryParamsMap["offset"], "5")
}
