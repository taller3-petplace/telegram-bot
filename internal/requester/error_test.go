package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type testErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (testError testErrorResponse) GetMessage() string {
	return testError.Message
}

func (testError testErrorResponse) GetStatus() int {
	return testError.Status
}

func TestErrorPolicyFunction(t *testing.T) {
	expectedErrorResponse := testErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "error te estas portando mal: seras castigada",
	}
	errorResponseRaw, err := json.Marshal(expectedErrorResponse)
	require.NoError(t, err)

	testCases := []struct {
		Name          string
		Response      *http.Response
		ExpectedError error
	}{
		{
			Name:          "If status code is less than 400, none error is returned",
			Response:      &http.Response{StatusCode: http.StatusOK},
			ExpectedError: nil,
		},
		{
			Name: "Error unmarshalling response",
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString("{wrong JSON format")),
			},
			ExpectedError: errUnmarshallingErrorResponse,
		},
		{
			Name: "Error response unmarshalled correctly",
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBuffer(errorResponseRaw)),
			},
			ExpectedError: fmt.Errorf("error te estas portando mal: seras castigada"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			errorResponse := ErrPolicyFunc[testErrorResponse](testCase.Response)
			if testCase.ExpectedError != nil {
				assert.ErrorContains(t, errorResponse, testCase.ExpectedError.Error())
				return
			}

			assert.Nil(t, errorResponse)
		})
	}
}

func TestRequestError(t *testing.T) {
	errorMessage := fmt.Errorf("hola perdida")
	t.Run("No content error", func(t *testing.T) {
		err := requestError{
			err:        errorMessage,
			statusCode: http.StatusNoContent,
			extraData:  "",
		}
		assert.True(t, err.IsNoContent())
		assert.False(t, err.IsBadRequest())
		assert.False(t, err.IsNotFound())
	})

	t.Run("Bad request error", func(t *testing.T) {
		err := requestError{
			err:        errorMessage,
			statusCode: http.StatusBadRequest,
			extraData:  "",
		}
		assert.True(t, err.IsBadRequest())
		assert.False(t, err.IsNoContent())
		assert.False(t, err.IsNotFound())
	})

	t.Run("Not found error", func(t *testing.T) {
		err := requestError{
			err:        errorMessage,
			statusCode: http.StatusNotFound,
			extraData:  "",
		}
		assert.True(t, err.IsNotFound())
		assert.False(t, err.IsNoContent())
		assert.False(t, err.IsBadRequest())
	})
}

func TestRequestError_Error(t *testing.T) {
	errorMessage := fmt.Errorf("te pones loquita de noche")
	extraData := "Olha gatinha, você me faz enlouquecer"
	t.Run("Error with extra data", func(t *testing.T) {
		err := requestError{
			err:        errorMessage,
			statusCode: http.StatusInternalServerError,
			extraData:  extraData,
		}

		expectedErrorMessage := fmt.Sprintf("500 - %v: %s", errorMessage, extraData)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})

	t.Run("Error without extra data", func(t *testing.T) {
		err := requestError{
			err:        errorMessage,
			statusCode: http.StatusInternalServerError,
			extraData:  "",
		}

		expectedErrorMessage := fmt.Sprintf("500 - %v", errorMessage)
		assert.Equal(t, expectedErrorMessage, err.Error())
	})
}
