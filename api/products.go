package api

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ProductsAPI provides operations to retrieve, create and update products.
type ProductsAPI struct {
	svc ProductsService
	r   *chi.Mux
}

type ProductsService interface {
	CreateProduct(categorySlug string, jsonData []byte) error
}

// NewProductsAPI returns a new ProductsAPI.
func newProductsAPI(svc ProductsService) http.Handler {
	a := &ProductsAPI{
		svc: svc,
		r:   chi.NewRouter(),
	}

	a.r.Post("/{category_slug}", a.CreateProduct)

	return a.r
}

func (a *ProductsAPI) CreateProduct(w http.ResponseWriter, r *http.Request) {
	categorySlug := chi.URLParam(r, "category_slug")

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	if err := a.svc.CreateProduct(categorySlug, jsonData); err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]any{
		"success": true,
	})
}
