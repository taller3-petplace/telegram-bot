package requester

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/utils/urlutils"
)

const getUser = "get_user"

// GetUserData fetches information about a user based on the given telegramID
func (r *Requester) GetUserData(telegramID int64) (domain.UserInfo, error) {
	operation := "GetUserData"
	endpointData, err := r.UsersService.GetEndpoint(getUser)
	if err != nil {
		logrus.Errorf("%v", err)
		return domain.UserInfo{}, err
	}

	url := endpointData.GetURL()
	url = urlutils.FormatURL(url, map[string]string{"telegramID": fmt.Sprintf("%v", telegramID)})
	request, err := http.NewRequest(endpointData.Method, url, nil)
	if err != nil {
		err = fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
		logrus.Errorf("%v", err)
		return domain.UserInfo{}, err
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing GetUserData: %v", err)
		return domain.UserInfo{}, NewRequestError(
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
		logrus.Error("nil response from users service")
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return domain.UserInfo{}, errorResponse
	}

	err = ErrPolicyFunc[petServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("error from users service: %v", err)
		return domain.UserInfo{}, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading users body: %v", err)
		return domain.UserInfo{}, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var userServiceResponse domain.UserServiceResponse
	err = json.Unmarshal(responseBody, &userServiceResponse)
	if err != nil {
		logrus.Errorf("error unmarshalling user data: %v", err)
		return domain.UserInfo{}, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingUserData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return userServiceResponse.UserData, nil
}
