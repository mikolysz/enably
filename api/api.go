package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mikolysz/enably/model"
)

type api struct {
	r *chi.Mux
}

// Dependencies contains all the dependencies that the API needs.
type Dependencies struct {
	Metadata         MetadataService
	Products         ProductsService
	Auth             AuthService
	ModerationAPIKey string
}

// New returns an http.Handler that responds to API requests.
func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	auth := newAuthAPI(deps.Auth)
	r.Use(auth.addAuthInfoToContext)

	r.Mount("/auth", auth)

	cat := newCategoriesAPI(deps.Metadata)
	r.Mount("/categories", cat)

	prod := newProductsAPI(deps.Products)
	r.Mount("/products", prod)

	mod := newModerationAPI(deps.Products, deps.ModerationAPIKey)
	r.Mount("/moderation", mod)

	return &api{r}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.r.ServeHTTP(w, r)
}

func jsonResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error when encoding response as JSON: %s", err)
	}
}

func errorResponse(w http.ResponseWriter, err error) {
	var uf model.UserFacingError
	if !errors.As(err, &uf) {
		uf = model.NewInternalServerError(err)
	}
	log.Printf("Error: %s", err)
	jsonResponse(w, uf.HTTPStatusCode, map[string]any{
		"type":    "error",
		"code":    uf.HTTPStatusCode,
		"message": uf.UserFacingMessage,
	})
}
