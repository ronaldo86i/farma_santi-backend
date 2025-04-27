package domain

type TokenResponse struct {
	Message      string `json:"message"`
	TokenAccess  string `json:"-"`
	TokenRefresh string `json:"-"`
}
