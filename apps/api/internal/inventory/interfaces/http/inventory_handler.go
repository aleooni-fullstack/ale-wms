package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type InventoryHandler struct {
	service *application.InventoryService
}

func NewInventoryHandler(service *application.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) RegisterRoutes(r chi.Router) {
	r.Route("/inventory", func(r chi.Router) {
		r.Post("/movements", h.registerMovement)
		r.Get("/movements", h.listMovements)
		r.Get("/balances", h.getBalance)
		r.Get("/balances/product/{productID}", h.listBalancesByProduct)
		r.Get("/balances/location/{locationID}", h.listBalancesByLocation)
	})
}

type registerMovementRequest struct {
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Type       string  `json:"type"`
	Quantity   float64 `json:"quantity"`
	Note       string  `json:"note"`
}

type movementResponse struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Type       string  `json:"type"`
	Quantity   float64 `json:"quantity"`
	Note       string  `json:"note"`
	CreatedAt  string  `json:"created_at"`
}

type balanceResponse struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	UpdatedAt  string  `json:"updated_at"`
}

type listMovementsResponse struct {
	Data    []movementResponse `json:"data"`
	Page    int32              `json:"page"`
	PerPage int32              `json:"per_page"`
}

func (h *InventoryHandler) registerMovement(w http.ResponseWriter, r *http.Request) {
	var req registerMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	movement, err := h.service.RegisterMovement(r.Context(), application.RegisterMovementInput{
		ProductID:  req.ProductID,
		LocationID: req.LocationID,
		Type:       req.Type,
		Quantity:   req.Quantity,
		Note:       req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toMovementResponse(movement))
}

func (h *InventoryHandler) listMovements(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.ListMovements(r.Context(), application.ListMovementsInput{
		ProductID:  r.URL.Query().Get("product_id"),
		LocationID: r.URL.Query().Get("location_id"),
		Page:       page,
		PerPage:    perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]movementResponse, len(result.Data))
	for i, m := range result.Data {
		resp[i] = toMovementResponse(m)
	}

	writeJSON(w, http.StatusOK, listMovementsResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *InventoryHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	locationID := r.URL.Query().Get("location_id")

	if productID == "" || locationID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "product_id and location_id are required"})
		return
	}

	balance, err := h.service.GetBalance(r.Context(), productID, locationID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toBalanceResponse(balance))
}

func (h *InventoryHandler) listBalancesByProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")

	balances, err := h.service.ListBalancesByProduct(r.Context(), productID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]balanceResponse, len(balances))
	for i, b := range balances {
		resp[i] = toBalanceResponse(b)
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *InventoryHandler) listBalancesByLocation(w http.ResponseWriter, r *http.Request) {
	locationID := chi.URLParam(r, "locationID")

	balances, err := h.service.ListBalancesByLocation(r.Context(), locationID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]balanceResponse, len(balances))
	for i, b := range balances {
		resp[i] = toBalanceResponse(b)
	}

	writeJSON(w, http.StatusOK, resp)
}

func toMovementResponse(m *domain.StockMovement) movementResponse {
	return movementResponse{
		ID:         m.ID,
		ProductID:  m.ProductID,
		LocationID: m.LocationID,
		Type:       string(m.Type),
		Quantity:   m.Quantity,
		Note:       m.Note,
		CreatedAt:  m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toBalanceResponse(b *domain.StockBalance) balanceResponse {
	return balanceResponse{
		ID:         b.ID,
		ProductID:  b.ProductID,
		LocationID: b.LocationID,
		Quantity:   b.Quantity,
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
	log.Println("ERROR:", err)

	var domErr *dderr.DomainError
	if errors.As(err, &domErr) {
		switch domErr.Code {
		case "INSUFFICIENT_STOCK", "INVALID_FILTER", "INVALID_PRODUCT_ID", "INVALID_LOCATION_ID", "INVALID_QUANTITY", "INVALID_MOVEMENT_TYPE":
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
