package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably/model"
)

type MetadataService interface {
	GetTopLevelCategories() []*model.Category
	GetCategory(slug string) (*model.Category, error)
}

type categoriesAPI struct {
	*chi.Mux
	meta MetadataService
}

func newCategoriesAPI(metadata MetadataService) *categoriesAPI {
	r := chi.NewRouter()
	c := &categoriesAPI{r, metadata}

	r.Get("/", c.getTopLevelCategories)
	r.Get("/{slug}", c.GetCategory)

	return c
}

// GetTopLevelCategories returns all categories that have no parent.
func (c *categoriesAPI) getTopLevelCategories(w http.ResponseWriter, r *http.Request) {
	cats := c.meta.GetTopLevelCategories()

	jsonResponse(w, http.StatusOK, cats)
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
