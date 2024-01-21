package domain

type UserInfo struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	City     string `json:"city"`
}
