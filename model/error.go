package model

import "net/http"

// UserFacingError is an error that can safely be shown to the user.
type UserFacingError struct {
	HTTPStatusCode    int    // The code to respond with when encountering this error
	UserFacingMessage string // The message to show to the user
	SecretMessage     string // The message to show in the logs, if any.
}

func (e UserFacingError) Error() string {
	if e.SecretMessage != "" {
		return e.SecretMessage
	}

	return e.UserFacingMessage
}

// NewInternalServerError returns a UserFacingError that indicates an internal server error.
func NewInternalServerError(err error) UserFacingError {
	return UserFacingError{
		HTTPStatusCode:    http.StatusInternalServerError,
		UserFacingMessage: "Internal server error",
		SecretMessage:     err.Error(),
	}
}
