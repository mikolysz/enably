package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably/model"
)

type MetadataService interface {
	GetRootCategory() *model.Category
	GetCategory(slug string) (*model.Category, error)
	GetFieldsetsForCategory(slug string) ([]*model.Fieldset, error)
	GetSchemaForFieldset(fieldset *model.Fieldset) (map[string]any, error)
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
	r.Get("/{slug}/fieldsets", c.GetFieldsetsForCategory)

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

type fieldset struct {
	Slug       string         `json:"slug"`
	Name       string         `json:"name"`
	Fields     []*model.Field `json:"fields"`
	JSONSchema map[string]any `json:"json_schema"`
}

func (C *categoriesAPI) GetFieldsetsForCategory(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	fsets, err := C.meta.GetFieldsetsForCategory(slug)
	if err != nil {
		errorResponse(w, err)
		return
	}

	// Get the JSON schems for these fieldsets.
	var withSchemas []*fieldset
	for _, fset := range fsets {
		schema, err := C.meta.GetSchemaForFieldset(fset)
		if err != nil {
			errorResponse(w, err)
			return
		}

		withSchemas = append(withSchemas, &fieldset{
			Slug:       fset.Slug,
			Name:       fset.Name,
			Fields:     fset.Fields,
			JSONSchema: schema,
		})
	}

	jsonResponse(w, http.StatusOK, withSchemas)
}
