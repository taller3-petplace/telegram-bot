package urlutils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestFormatURL(t *testing.T) {
	testCases := []struct {
		Name        string
		InputURL    string
		Params      map[string]string
		ExpectedURL string
	}{
		{
			Name:     "Replace multiple params",
			InputURL: "breaking-bad/{season}/{chapter}/best/episode/ever",
			Params: map[string]string{
				"season":  "S5",
				"chapter": "Ch14",
			},
			ExpectedURL: "breaking-bad/S5/Ch14/best/episode/ever",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			formattedURL := FormatURL(testCase.InputURL, testCase.Params)
			assert.Equal(t, testCase.ExpectedURL, formattedURL)
		})
	}
}

func TestSarasa(t *testing.T) {
	url := "https://te/estas/portando/mal/{owner_id}"

	url = FormatURL(url, map[string]string{"owner_id": "12345"})

	request, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	queryParamsValues := request.URL.Query()
	queryParamsValues.Add("limit", "100")
	queryParamsValues.Add("offset", "0")

	request.URL.RawQuery = queryParamsValues.Encode()

	fmt.Printf("%s", request.URL.String())
}
