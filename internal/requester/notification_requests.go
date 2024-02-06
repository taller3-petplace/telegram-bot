package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"telegram-bot/internal/domain"
)

const scheduleNotifications = "schedule_notifications"

// RegisterNotifications sends a request to Notification Scheduler service to create multiple notifications
// with the provided data in domain.NotificationRequest
func (r *Requester) RegisterNotifications(notificationRequest domain.NotificationRequest) ([]domain.NotificationResponse, error) {
	operation := "ScheduleNotifications"
	endpointData, err := r.NotificationsService.GetEndpoint(scheduleNotifications)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}

	url := endpointData.GetURL()
	rawBody, err := json.Marshal(notificationRequest)
	if err != nil {
		logrus.Errorf("error marshalling notification request: %v", err)
		return nil, fmt.Errorf("%w: %v", errMarshallingNotificationRequest, err)
	}

	requestBody := bytes.NewReader(rawBody)
	request, err := http.NewRequest(endpointData.Method, url, requestBody)
	if err != nil {
		err = fmt.Errorf("%w: %v. Operation: %s", errCreatingRequest, err, operation)
		logrus.Errorf("%v", err)
		return nil, err
	}

	setTelegramHeader(request)
	response, err := r.clientHTTP.Do(request)
	if err != nil {
		logrus.Errorf("error performing ScheduleNotifications: %v", err)
		return nil, NewRequestError(
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
		logrus.Error("nil response from notifications service")
		errorResponse := NewRequestError(
			errNilResponse,
			http.StatusInternalServerError,
			operation,
		)
		return nil, errorResponse
	}

	err = ErrPolicyFunc[notificationServiceErrorResponse](response)
	if err != nil {
		logrus.Errorf("error from notifications service: %v", err)
		return nil, NewRequestError(
			err,
			response.StatusCode,
			"",
		)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("error reading notification body: %v", err)
		return nil, NewRequestError(
			errReadingResponseBody,
			http.StatusInternalServerError,
			operation,
		)
	}

	var notificationsResponse []domain.NotificationResponse
	err = json.Unmarshal(responseBody, &notificationsResponse)
	if err != nil {
		logrus.Errorf("error unmarshalling notification data: %v", err)
		return nil, NewRequestError(
			fmt.Errorf("%w: %v", errUnmarshallingNotificationsData, err),
			http.StatusInternalServerError,
			"",
		)
	}

	return notificationsResponse, nil
}
