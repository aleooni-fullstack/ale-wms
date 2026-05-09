package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type WarehouseHandler struct {
	service *application.WarehouseService
}

func NewWarehouseHandler(service *application.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) RegisterRoutes(r chi.Router) {
	r.Route("/warehouses", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.deactivate)
	})
}

type createWarehouseRequest struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type updateWarehouseRequest struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type warehouseResponse struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Active  bool   `json:"active"`
}

type listWarehousesResponse struct {
	Data    []warehouseResponse `json:"data"`
	Page    int32               `json:"page"`
	PerPage int32               `json:"per_page"`
}

func (h *WarehouseHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListWarehousesInput{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]warehouseResponse, len(result.Data))
	for i, wh := range result.Data {
		resp[i] = toWarehouseResponse(wh)
	}

	writeJSON(w, http.StatusOK, listWarehousesResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *WarehouseHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	warehouse, err := h.service.Create(r.Context(), application.CreateWarehouseInput{
		Code:    req.Code,
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toWarehouseResponse(warehouse))
}

func (h *WarehouseHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	warehouse, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toWarehouseResponse(warehouse))
}

func (h *WarehouseHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateWarehouseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	warehouse, err := h.service.Update(r.Context(), id, application.UpdateWarehouseInput{
		Code:    req.Code,
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toWarehouseResponse(warehouse))
}

func (h *WarehouseHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Deactivate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toWarehouseResponse(w *domain.Warehouse) warehouseResponse {
	return warehouseResponse{
		ID:      w.ID,
		Code:    w.Code,
		Name:    w.Name,
		Address: w.Address,
		Active:  w.Active,
	}
}

func parsePagination(r *http.Request) (page, perPage int32) {
	page = 1
	perPage = 20

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

	return page, perPage
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
		case "CODE_ALREADY_EXISTS", "INVALID_CODE", "INVALID_NAME", "INVALID_WAREHOUSE_ID", "INVALID_ZONE_ID":
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
