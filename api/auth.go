package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authApi struct {
	svc AuthService
	r   *chi.Mux
}

type AuthService interface {
	SendLoginEmail(c context.Context, email, redirectURI string) error
	AuthenticateUser(c context.Context, token string) (email string, err error)
}

func newAuthAPI(svc AuthService) *authApi {
	a := &authApi{
		svc: svc,
		r:   chi.NewRouter(),
	}

	a.r.Post("/login", a.SendLoginEmail)
	a.r.Get("/me", a.GetUserInfo)
	return a
}

func (a *authApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func (a *authApi) SendLoginEmail(w http.ResponseWriter, r *http.Request) {
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	type Data struct {
		Email       string `json:"email"`
		RedirectURI string `json:"redirect_uri"`
	}

	var data Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		errorResponse(w, err)
		return
	}

	if err := a.svc.SendLoginEmail(r.Context(), data.Email, data.RedirectURI); err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (a *authApi) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// middleware

type emailContextKey string

func (a *authApi) addAuthInfoToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Header := r.Header.Get("Authorization")
		if len(Header) < 8 || Header[:7] != "Bearer " {
			next.ServeHTTP(w, r)
			return
		}

		token := Header[7:]
		email, err := a.svc.AuthenticateUser(r.Context(), token)
		if err != nil {
			errorResponse(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), emailContextKey("email"), email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
