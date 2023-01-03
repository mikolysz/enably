package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably/model"
)

type MetadataService interface {
	GetRootCategory() *model.Category
	GetCategory(slug string) (*model.Category, error)
	GetSchemasForCategory(category *model.Category) (map[string]any, error)
}

type categoriesAPI struct {
	*chi.Mux
	meta MetadataService
}

func newCategoriesAPI(metadata MetadataService) *categoriesAPI {
	r := chi.NewRouter()
	c := &categoriesAPI{r, metadata}

	r.Get("/", c.getRootCategory)
	r.Get("/{slug}", c.GetCategory)
	r.Get("/{slug}/schemas", c.GetSchemasForCategory)

	return c
}

func (c *categoriesAPI) getRootCategory(w http.ResponseWriter, r *http.Request) {
	cat := c.meta.GetRootCategory()

	jsonResponse(w, http.StatusOK, cat)
}

func (c *categoriesAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cat, err := c.meta.GetCategory(slug)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, cat)
}

func (c *categoriesAPI) GetSchemasForCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cat, err := c.meta.GetCategory(slug)
	if err != nil {
		errorResponse(w, err)
		return
	}
	schemas, err := c.meta.GetSchemasForCategory(cat)
	if err != nil {
		errorResponse(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, schemas)
}
