package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/catalog/domain"
)

type ProductHandler struct {
	service *application.ProductService
}

type listProductsResponse struct {
	Data    []productResponse `json:"data"`
	Page    int32             `json:"page"`
	PerPage int32             `json:"per_page"`
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.deactivate)
	})
}

type createProductRequest struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}

type updateProductRequest struct {
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}

type productResponse struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	Active      bool   `json:"active"`
}

func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	page := int32(1)
	perPage := int32(20)

	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = int32(n)
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			perPage = int32(n)
		}
	}

	result, err := h.service.List(r.Context(), application.ListProductsInput{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]productResponse, len(result.Data))
	for i, p := range result.Data {
		resp[i] = toResponse(p)
	}

	writeJSON(w, http.StatusOK, listProductsResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	product, err := h.service.Create(r.Context(), application.CreateProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toResponse(product))
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(product))
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	product, err := h.service.Update(r.Context(), id, application.UpdateProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Unit:        req.Unit,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(product))
}

func (h *ProductHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Deactivate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toResponse(p *domain.Product) productResponse {
	return productResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Unit:        p.Unit,
		Active:      p.Active,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	var domErr *dderr.DomainError
	if errors.As(err, &domErr) {
		switch domErr.Code {
		case "SKU_ALREADY_EXISTS", "INVALID_SKU", "INVALID_NAME", "INVALID_UNIT":
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": domErr.Message})
			return
		}
	}

	if errors.Is(err, dderr.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
