package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type InventoryBalanceHandler struct {
	service *application.InventoryBalanceService
}

func NewInventoryBalanceHandler(service *application.InventoryBalanceService) *InventoryBalanceHandler {
	return &InventoryBalanceHandler{service: service}
}

func (h *InventoryBalanceHandler) RegisterRoutes(r chi.Router) {
	r.Route("/inventory/balances/count", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
		r.Post("/{id}/items", h.addItem)
		r.Put("/{id}/items/{productID}/count", h.countItem)
	})
}

type createInventoryBalanceRequest struct {
	LocationID string `json:"location_id"`
	Note       string `json:"note"`
}

type addItemRequest struct {
	ProductID string `json:"product_id"`
}

type countItemRequest struct {
	CountedQuantity float64 `json:"counted_quantity"`
}

type inventoryBalanceItemResponse struct {
	ID                 string   `json:"id"`
	InventoryBalanceID string   `json:"inventory_balance_id"`
	ProductID          string   `json:"product_id"`
	SystemQuantity     float64  `json:"system_quantity"`
	CountedQuantity    *float64 `json:"counted_quantity"`
	Difference         *float64 `json:"difference"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

type inventoryBalanceResponse struct {
	ID         string                         `json:"id"`
	LocationID string                         `json:"location_id"`
	Status     string                         `json:"status"`
	Note       string                         `json:"note"`
	Items      []inventoryBalanceItemResponse `json:"items"`
	CreatedAt  string                         `json:"created_at"`
	UpdatedAt  string                         `json:"updated_at"`
}

type listInventoryBalancesResponse struct {
	Data    []inventoryBalanceResponse `json:"data"`
	Page    int32                      `json:"page"`
	PerPage int32                      `json:"per_page"`
}

func (h *InventoryBalanceHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListInventoryBalancesInput{
		LocationID: r.URL.Query().Get("location_id"),
		Page:       page,
		PerPage:    perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]inventoryBalanceResponse, len(result.Data))
	for i, b := range result.Data {
		resp[i] = toInventoryBalanceResponse(b)
	}

	writeJSON(w, http.StatusOK, listInventoryBalancesResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *InventoryBalanceHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createInventoryBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	balance, err := h.service.Create(r.Context(), application.CreateInventoryBalanceInput{
		LocationID: req.LocationID,
		Note:       req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toInventoryBalanceResponse(balance))
}

func (h *InventoryBalanceHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	balance, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toInventoryBalanceResponse(balance))
}

func (h *InventoryBalanceHandler) start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	balance, err := h.service.Start(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toInventoryBalanceResponse(balance))
}

func (h *InventoryBalanceHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	balance, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toInventoryBalanceResponse(balance))
}

func (h *InventoryBalanceHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	balance, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toInventoryBalanceResponse(balance))
}

func (h *InventoryBalanceHandler) addItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.AddItem(r.Context(), application.AddItemInput{
		InventoryBalanceID: id,
		ProductID:          req.ProductID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toInventoryBalanceItemResponse(item))
}

func (h *InventoryBalanceHandler) countItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID := chi.URLParam(r, "productID")

	var req countItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.CountItem(r.Context(), application.CountItemInput{
		InventoryBalanceID: id,
		ProductID:          productID,
		CountedQuantity:    req.CountedQuantity,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toInventoryBalanceItemResponse(item))
}

func toInventoryBalanceResponse(b *domain.InventoryBalance) inventoryBalanceResponse {
	items := make([]inventoryBalanceItemResponse, len(b.Items))
	for i, item := range b.Items {
		items[i] = toInventoryBalanceItemResponse(item)
	}

	return inventoryBalanceResponse{
		ID:         b.ID,
		LocationID: b.LocationID,
		Status:     string(b.Status),
		Note:       b.Note,
		Items:      items,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toInventoryBalanceItemResponse(item *domain.InventoryBalanceItem) inventoryBalanceItemResponse {
	return inventoryBalanceItemResponse{
		ID:                 item.ID,
		InventoryBalanceID: item.InventoryBalanceID,
		ProductID:          item.ProductID,
		SystemQuantity:     item.SystemQuantity,
		CountedQuantity:    item.CountedQuantity,
		Difference:         item.Difference(),
		CreatedAt:          item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
