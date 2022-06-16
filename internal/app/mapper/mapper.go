package mapper

import "fmt"

type LoginRequest struct {
	UserID string `json:"user_id"`
}

func (l LoginRequest) Bind() error {
	if len(l.UserID) == 0 {
		return fmt.Errorf("user_id must not be empty")
	}
	return nil
}

type TokenPair struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (r TokenPair) Bind() error {
	if len(r.AccessToken) == 0 || len(r.RefreshToken) == 0 {
		return fmt.Errorf("access_token and refresh_token must not be empty")
	}
	return nil
}
