package app

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mikolysz/enably/model"
	"github.com/mikolysz/enably/pkg/email"
)

// AuthenticationService manages user authentication
type AuthenticationService struct {
	store       TokenStore
	emailSender email.Sender
	frontendURL *url.URL
}

type TokenStore interface {
	AddToken(c context.Context, t model.Token) (model.Token, error)
	GetByTokenContents(c context.Context, token string) (model.Token, error)
}

// NewAuthenticationService creates a new AuthenticationService.
func NewAuthenticationService(store TokenStore, emailSender email.Sender, frontendURL *url.URL) AuthenticationService {
	return AuthenticationService{
		store:       store,
		emailSender: emailSender,
		frontendURL: frontendURL,
	}
}

// SendLoginEmail creates an authentication token and emails the user with a "magic link" which logs them in.
// redirectURI is the URL to redirect to after the login is successful.
func (s AuthenticationService) SendLoginEmail(c context.Context, emailAddress, redirectURI string) error {
	token := model.NewToken(emailAddress)
	if _, err := s.store.AddToken(c, *token); err != nil {
		return err
	}

	authValues := url.Values{}
	authValues.Set("token", token.Token)
	// FIXME: Use different tokens for initial authentication and keeping the session, that's more secure.
	authValues.Set("redirect_uri", redirectURI)

	url := *s.frontendURL
	url.Path = "authorize"
	url.RawQuery = authValues.Encode()
	urlStr := url.String()
	msg := email.Message{
		Recipient:        emailAddress,
		Subject:          "Login to Enably",
		PlainTextContent: "Click this link to login: " + urlStr,
		HTMLContent:      "<p>Click this link to login: <a href=\"" + urlStr + "\">Login</a></p>",
	}

	if err := s.emailSender.Send(msg); err != nil {
		return fmt.Errorf("error when sending login: %w", err)
	}
	// FIXME: security: ensure redirectURI points to our own domain just in case.

	return nil
}

// AuthenticateUser checks if the token is valid and returns the email address of the user.
func (s AuthenticationService) AuthenticateUser(c context.Context, token string) (email string, err error) {
	t, err := s.store.GetByTokenContents(c, token)
	if err != nil {
		return "", err
	}
	return t.Email, nil
}
