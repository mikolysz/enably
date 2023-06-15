package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mikolysz/enably/model"
)

type moderationAPI struct {
	svc    ProductsService
	apiKey string
	r      *chi.Mux
}

// NewModerationAPI returns a new ModerationAPI.
func newModerationAPI(svc ProductsService, apiKey string) http.Handler {
	a := &moderationAPI{
		svc:    svc,
		apiKey: apiKey,
		r:      chi.NewRouter(),
	}

	a.r.Use(a.checkAPIKey)
	a.r.Get("/pending", a.GetPendingProducts)
	a.r.Post("/products/{product_id}/approve", a.ApproveProduct)
	a.r.Post("/products/{product_id}/reject", a.RejectProduct)
	return a.r
}

// checkAPIKey is a middleware that returns 403 if the moderation API key is wrong or missing.
func (a *moderationAPI) checkAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Moderation-API-Key")
		if key == "" {
			errorResponse(w, model.UserFacingError{
				HTTPStatusCode:    http.StatusForbidden,
				UserFacingMessage: "missing moderation API key.\n\nPlease provide the X-Moderation-API-Key header",
			})
			return
		}
		if key != a.apiKey {
			errorResponse(w, model.UserFacingError{
				HTTPStatusCode:    http.StatusForbidden,
				UserFacingMessage: "wrong moderation API key",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *moderationAPI) GetPendingProducts(w http.ResponseWriter, r *http.Request) {
	prods, err := a.svc.GetProductsNeedingApproval()
	if err != nil {
		errorResponse(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, prods)
}

func (a *moderationAPI) ApproveProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, model.UserFacingError{
			HTTPStatusCode:    http.StatusBadRequest,
			UserFacingMessage: "invalid product ID",
		})
		return
	}

	if err := a.svc.ApproveProduct(id); err != nil {
		errorResponse(w, err)
		return
	}
}

func (a *moderationAPI) RejectProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, model.UserFacingError{
			HTTPStatusCode:    http.StatusBadRequest,
			UserFacingMessage: "invalid product ID",
		})
		return
	}

	if err := a.svc.RejectProduct(id); err != nil {
		errorResponse(w, err)
		return
	}
}
