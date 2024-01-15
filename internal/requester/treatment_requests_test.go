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

func TestRequesterGetTreatmentsByPetID(t *testing.T) {
	currentTime := time.Now()
	treatmentsServiceEndpoints := getExpectedTreatmentsServiceEndpoints()
	getPetTreatmentsEndpoint := treatmentsServiceEndpoints[getPetTreatments]
	getPetTreatmentsEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := treatmentsServiceEndpoints[getPets]
	invalidEndpoint.Method = "hola perdida"

	requester := Requester{
		TreatmentsService: config.ServiceEndpoints{
			Endpoints: treatmentsServiceEndpoints,
		},
	}

	treatmentsServiceError := treatmentServiceErrorResponse{
		Code: http.StatusInternalServerError,
		Msg:  "error ",
	}
	serviceErrorRaw, err := json.Marshal(treatmentsServiceError)
	require.NoError(t, err)

	oldestTreatment := domain.Treatment{
		ID:           "1",
		Type:         "Medical appointment",
		Comments:     []domain.Comment{},
		DateStart:    currentTime.AddDate(-1, 0, 0),
		LastModified: currentTime.AddDate(-1, 0, 0),
	}

	newestTreatment := domain.Treatment{
		ID:           "2",
		Type:         "Surgery",
		Comments:     []domain.Comment{},
		DateStart:    currentTime.AddDate(0, -1, 0),
		LastModified: currentTime.AddDate(0, -1, 0),
	}

	treatmentsData := []domain.Treatment{
		oldestTreatment,
		newestTreatment,
	}

	expectedTreatments := []domain.Treatment{
		newestTreatment,
		oldestTreatment,
	}

	rawTreatmentsData, err := json.Marshal(treatmentsData)
	require.NoError(t, err)

	testCases := []struct {
		Name               string
		Requester          Requester
		ClientMockConfig   *clientMockConfig
		ExpectsError       bool
		ExpectedError      error
		ExpectedTreatments []domain.Treatment
	}{
		{
			Name: "Endpoint does not exist",
			Requester: Requester{
				TreatmentsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{}},
			},
			ExpectsError:  true,
			ExpectedError: errEndpointDoesNotExist,
		},
		{
			Name: "Error creating request",
			Requester: Requester{
				TreatmentsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{
					getPetTreatments: invalidEndpoint,
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
			Name:      "Error from treatments service",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(serviceErrorRaw)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: fmt.Errorf(treatmentsServiceError.GetMessage()),
		},
		{
			Name:      "Error unmarshalling treatments data",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id": "69abc"}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingMultipleTreatments,
		},
		{
			Name:      "Get treatments correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawTreatmentsData)),
				},
				Err: nil,
			},
			ExpectsError:       false,
			ExpectedTreatments: expectedTreatments,
			ExpectedError:      nil,
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

			treatmentsResponse, err := testCase.Requester.GetTreatmentsByPetID(petID)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedTreatments, treatmentsResponse)
		})
	}
}

func TestRequesterGetTreatment(t *testing.T) {
	currentTime := time.Now()
	treatmentsServiceEndpoints := getExpectedTreatmentsServiceEndpoints()
	getTreatmentEndpoint := treatmentsServiceEndpoints[getTreatment]
	getTreatmentEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := treatmentsServiceEndpoints[getTreatment]
	invalidEndpoint.Method = "japiaguar - me vas a extrañar"

	requester := Requester{
		TreatmentsService: config.ServiceEndpoints{
			Endpoints: treatmentsServiceEndpoints,
		},
	}

	treatmentsServiceError := treatmentServiceErrorResponse{
		Code: http.StatusInternalServerError,
		Msg:  "error te pones loquita de noche",
	}

	serviceErrorRaw, err := json.Marshal(treatmentsServiceError)
	require.NoError(t, err)

	treatmentID := 1
	// to check that comments are sorted
	oldestComment := domain.Comment{
		DateAdded:   currentTime.AddDate(0, -1, 0),
		Information: "Lloraras mas de diez veces por amor",
		Owner:       "Leo Mattioli",
	}
	newestComment := domain.Comment{
		DateAdded:   currentTime,
		Information: "quiero una chica quiero una ya, quiero una mujer que sea muy especial",
		Owner:       "Escucha",
	}
	treatmentData := domain.Treatment{
		ID:   fmt.Sprint(treatmentID),
		Type: "Medical appointment",
		Comments: []domain.Comment{
			oldestComment,
			newestComment,
		},
		DateStart:    currentTime.AddDate(0, -1, 0),
		LastModified: currentTime,
	}

	rawTreatmentData, err := json.Marshal(treatmentData)
	require.NoError(t, err)

	expectedTreatmentData := treatmentData
	expectedTreatmentData.Comments = []domain.Comment{newestComment, oldestComment}

	testCases := []struct {
		Name                  string
		Requester             Requester
		ClientMockConfig      *clientMockConfig
		ExpectsError          bool
		ExpectedError         error
		ExpectedTreatmentData domain.Treatment
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
				TreatmentsService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{
					getTreatment: invalidEndpoint,
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
			Name:      "Error from treatments service",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(serviceErrorRaw)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: fmt.Errorf(treatmentsServiceError.GetMessage()),
		},
		{
			Name:      "Error unmarshalling treatment data",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id": 69}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingTreatmentData,
		},
		{
			Name:      "Get treatment data correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawTreatmentData)),
				},
				Err: nil,
			},
			ExpectsError:          false,
			ExpectedTreatmentData: expectedTreatmentData,
			ExpectedError:         nil,
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
			treatmentResponse, err := testCase.Requester.GetTreatment(treatmentID)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)

			// We need to check one by one, otherwise we get an error
			testCase.ExpectedTreatmentData.DateStart = currentTime
			treatmentResponse.DateStart = currentTime

			assert.Equal(t, testCase.ExpectedTreatmentData, treatmentResponse, "treatment data do not match")
		})
	}
}

func TestRequesterGetVaccines(t *testing.T) {
	// ToDo: implement me. Licha
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
