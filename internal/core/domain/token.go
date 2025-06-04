package domain

import "time"

type TokenResponse struct {
	Message         string    `json:"message"`
	AccessToken     string    `json:"-"`
	RefreshToken    string    `json:"-"`
	ExpAccessToken  time.Time `json:"-"`
	ExpRefreshToken time.Time `json:"-"`
}
