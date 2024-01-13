package urlutils

import (
	"fmt"
	"net/http"
	"regexp"
)

// FormatURL formats the given URL setting the values from the map to it. Is not-in-place
func FormatURL(url string, params map[string]string) string {
	formattedURL := url
	for param, value := range params {
		regex := regexp.MustCompile(fmt.Sprintf("{%s}", param))
		formattedURL = regex.ReplaceAllString(formattedURL, value)
	}

	return formattedURL
}

// AddQueryParams adds the given query params to the request. Is in-place
func AddQueryParams(request *http.Request, queryParams map[string]string) {
	queryParamsValues := request.URL.Query()
	for queryParam, value := range queryParams {
		queryParamsValues.Add(queryParam, value)
	}
	request.URL.RawQuery = queryParamsValues.Encode()
}
