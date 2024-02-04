package sender

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
