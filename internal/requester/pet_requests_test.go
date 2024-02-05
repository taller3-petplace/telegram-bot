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
	petID       = 1
	testBaseURL = "https://test"
)

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

	petsDataResponse := domain.PetsResponse{
		PetsData: []domain.PetData{
			{
				PetDataIdentifier: domain.PetDataIdentifier{
					ID:   1,
					Name: "Cartucho",
					Type: "DOG",
				},
			},
			{
				PetDataIdentifier: domain.PetDataIdentifier{
					ID:   2,
					Name: "Pantufla",
					Type: "CAT",
				},
			},
		},
		Paging: domain.Paging{
			Total:  2,
			Offset: 0,
			Limit:  100,
		},
	}

	rawResponse, err := json.Marshal(petsDataResponse)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		Requester        Requester
		ClientMockConfig *clientMockConfig
		ExpectsError     bool
		ExpectedError    error
		ExpectedPetsData []domain.PetData
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
					Body:       io.NopCloser(bytes.NewBufferString(`{"results": [{"id": "69abc"}]}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingMultiplePetsData,
		},
		{
			Name:      "Get pets data correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawResponse)),
				},
				Err: nil,
			},
			ExpectsError:     false,
			ExpectedPetsData: petsDataResponse.PetsData,
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
		OwnerID:      fmt.Sprint(ownerID),
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

func TestRequesterGetPetDataByID(t *testing.T) {
	petsServiceEndpoints := getExpectedPetsServiceEndpoints()
	getPetsByIDEndpoint := petsServiceEndpoints[getPetByID]
	getPetsByIDEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := petsServiceEndpoints[getPetByID]
	invalidEndpoint.Method = "japiaguar - bruja"

	requester := Requester{
		PetsService: config.ServiceEndpoints{
			Endpoints: petsServiceEndpoints,
		},
	}

	petsServiceError := petServiceErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "error te pones loquita de noche",
	}
	serviceErrorRaw, err := json.Marshal(petsServiceError)
	require.NoError(t, err)

	petData := domain.PetData{
		PetDataIdentifier: domain.PetDataIdentifier{
			ID:   petID,
			Name: "Cartucho",
			Type: "DOG",
		},
		Race:      "Perro salchicha",
		BirthDate: time.Now().Truncate(0),
	}

	rawPetData, err := json.Marshal(petData)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		Requester        Requester
		ClientMockConfig *clientMockConfig
		ExpectsError     bool
		ExpectedError    error
		ExpectedPetData  domain.PetData
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
					getPetByID: invalidEndpoint,
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
			Name:      "Error unmarshalling pet data",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id": "69abc"}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingPetData,
		},
		{
			Name:      "Get pet data correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawPetData)),
				},
				Err: nil,
			},
			ExpectsError:    false,
			ExpectedPetData: petData,
			ExpectedError:   nil,
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
			petDataResponse, err := testCase.Requester.GetPetData(petID)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedPetData, petDataResponse, "pet data do not match")
		})
	}
}
