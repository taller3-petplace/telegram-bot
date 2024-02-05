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
)

const telegramID = 911

func TestRequesterGetUser(t *testing.T) {
	usersServiceEndpoints := getExpectedUsersServiceEndpoints()
	getUserEndpoint := usersServiceEndpoints[getUser]
	getUserEndpoint.SetBaseURL(testBaseURL)

	invalidEndpoint := usersServiceEndpoints[getPets]
	invalidEndpoint.Method = "Mi nombre es Fernando " +
		"Soy el rey del contrabando " +
		"Y esta noche voy a arrasar"

	requester := Requester{
		UsersService: config.ServiceEndpoints{
			Endpoints: usersServiceEndpoints,
		},
	}

	// ToDo: replace this error. Licha
	usersServiceError := petServiceErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "error cae el soooool en tu balcoooooon",
	}
	serviceErrorRaw, err := json.Marshal(usersServiceError)
	require.NoError(t, err)

	userServiceResponse := domain.UserServiceResponse{
		UserData: domain.UserInfo{
			UserID:   "69",
			FullName: "Larry Capija",
			Email:    "larrydick@testmail.com",
			City:     "комната твоей сестры",
		},
		Code: http.StatusOK,
	}

	rawResponse, err := json.Marshal(userServiceResponse)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		Requester        Requester
		ClientMockConfig *clientMockConfig
		ExpectsError     bool
		ExpectedError    error
		ExpectedUserData domain.UserInfo
	}{
		{
			Name: "Endpoint does not exist",
			Requester: Requester{
				UsersService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{}},
			},
			ExpectsError:  true,
			ExpectedError: errEndpointDoesNotExist,
		},
		{
			Name: "Error creating request",
			Requester: Requester{
				UsersService: config.ServiceEndpoints{Endpoints: map[string]config.Endpoint{
					getUser: invalidEndpoint,
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
			Name:      "Error from users service",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBuffer(serviceErrorRaw)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: fmt.Errorf(usersServiceError.GetMessage()),
		},
		{
			Name:      "Error unmarshalling user data",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"code": "69abc"}`)),
				},
				Err: nil,
			},
			ExpectsError:  true,
			ExpectedError: errUnmarshallingUserData,
		},
		{
			Name:      "Get user data correctly",
			Requester: requester,
			ClientMockConfig: &clientMockConfig{
				ResponseBody: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(rawResponse)),
				},
				Err: nil,
			},
			ExpectsError:     false,
			ExpectedUserData: userServiceResponse.UserData,
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

			petsDataResponse, err := testCase.Requester.GetUserData(telegramID)
			if testCase.ExpectsError {
				assert.ErrorContains(t, err, testCase.ExpectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedUserData, petsDataResponse)
		})
	}
}
