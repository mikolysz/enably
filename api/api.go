package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikolysz/enably/model"
)

type api struct {
	r *chi.Mux
}

// Dependencies contains all the dependencies that the API needs.
type Dependencies struct {
	Metadata MetadataService
}

// New returns an http.Handler that responds to API requests.
func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	cat := newCategoriesAPI(deps.Metadata)
	r.Mount("/categories", cat)

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
