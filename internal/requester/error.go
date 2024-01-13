package requester

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const nonErrorCodeLimit = 299

// Services error responses definitions

type petServiceErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (petError petServiceErrorResponse) GetMessage() string {
	return petError.Message
}

func (petError petServiceErrorResponse) GetStatus() int {
	return petError.Status
}

type serviceError interface {
	GetMessage() string
	GetStatus() int
}

var (
	errEndpointDoesNotExist  = errors.New("error endpoint does not exist")
	errPerformingRequest     = errors.New("error performing request")
	errReadingResponseBody   = errors.New("error reading response body")
	errUnmarshallingPetsData = errors.New("error unmarshalling pets data")
	errMarshallingPetRequest = errors.New("error marshalling pet request")
	errCreatingRequest       = errors.New("error creating request")
	errNilResponse           = errors.New("error nil response")
)

func ErrPolicyFunc[serviceErrorType serviceError](response *http.Response) error {
	if response.StatusCode < http.StatusBadRequest {
		return nil
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("%w: cannot read error response body: %v", errReadingResponseBody, err)
	}

	var errResponse serviceErrorType
	err = json.Unmarshal(responseBody, &errResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling error information")
	}

	return fmt.Errorf("%s", errResponse.GetMessage())
}

type RequestError interface {
	error
	IsNoContent() bool
	IsBadRequest() bool
	IsNotFound() bool

	StatusCode() int
}

type requestError struct {
	err        error
	statusCode int
	extraData  string
}

func NewRequestError(err error, statusCode int, extraData string) error {
	return requestError{
		err:        err,
		statusCode: statusCode,
		extraData:  extraData,
	}
}

func (re requestError) Error() string {
	if re.extraData == "" {
		return fmt.Sprintf("%d - %v", re.statusCode, re.err)
	}

	return fmt.Sprintf("%d - %v: %s", re.statusCode, re.err, re.extraData)
}

func (re requestError) Is(target error) bool {
	return errors.Is(re.err, target)
}

func (re requestError) IsNoContent() bool {
	return re.statusCode == http.StatusNoContent
}

func (re requestError) IsBadRequest() bool {
	return re.statusCode == http.StatusBadRequest
}

func (re requestError) IsNotFound() bool {
	return re.statusCode == http.StatusNotFound
}

func (re requestError) StatusCode() int {
	return re.statusCode
}
