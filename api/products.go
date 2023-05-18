package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably/model"
)

// ProductsAPI provides operations to retrieve, create and update products.
type ProductsAPI struct {
	svc ProductsService
	r   *chi.Mux
}

type ProductsService interface {
	CreateProduct(categorySlug string, jsonData []byte) (model.Product, error)
	GetProductsByCategory(categorySlug string) ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)
}

// NewProductsAPI returns a new ProductsAPI.
func newProductsAPI(svc ProductsService) http.Handler {
	a := &ProductsAPI{
		svc: svc,
		r:   chi.NewRouter(),
	}

	a.r.Post("/{category_slug}", a.CreateProduct)
	a.r.Get("/by-category/{category_slug}", a.GetProductsByCategory)
	a.r.Get("/{product_id}", a.GetProductByID)
	return a.r
}

func (a *ProductsAPI) CreateProduct(w http.ResponseWriter, r *http.Request) {
	categorySlug := chi.URLParam(r, "category_slug")

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, err)
		return
	}

	prod, err := a.svc.CreateProduct(categorySlug, jsonData)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, prod)
}

func (a *ProductsAPI) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categorySlug := chi.URLParam(r, "category_slug")

	prods, err := a.svc.GetProductsByCategory(categorySlug)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, prods)
}

func (a *ProductsAPI) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	// Convert the ID to int, return bad request if failed
	id, err := strconv.Atoi(productID)
	if err != nil {
		errorResponse(w, model.UserFacingError{
			HTTPStatusCode:    http.StatusBadRequest,
			UserFacingMessage: "Invalid product ID",
		})
	}

	product, err := a.svc.GetProductByID(id)
	if err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, product)
}
