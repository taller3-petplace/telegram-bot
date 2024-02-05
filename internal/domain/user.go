package domain

// UserServiceResponse json response from users service when data of a user is requested
type UserServiceResponse struct {
	UserData UserInfo `json:"data"`
	Code     int      `json:"code"`
}

type UserInfo struct {
	UserID   string `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	City     string `json:"city"`
}
