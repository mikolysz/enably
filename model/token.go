package model

import (
	"net/http"

	"github.com/google/uuid"
)

// Token is a secure token authenticating a user with a specific email address.
type Token struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"created_at"`
}

var ErrInvalidToken = UserFacingError{
	HTTPStatusCode:    http.StatusForbidden,
	UserFacingMessage: "Invalid authentication token",
}

// NewToken 		creates a new token for a given email address.
func NewToken(email string) *Token {
	return &Token{
		Email: email,
		Token: uuid.NewString(),
	}
}
