package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type OrderHandler struct {
	service *application.OrderService
}

func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/items", h.addItem)
		r.Post("/{id}/confirm", h.confirm)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type createOrderRequest struct {
	Reference string `json:"reference"`
	Note      string `json:"note"`
}

type addOrderItemRequest struct {
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
}

type orderItemResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"order_id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	CreatedAt  string  `json:"created_at"`
}

type orderResponse struct {
	ID        string              `json:"id"`
	Reference string              `json:"reference"`
	Status    string              `json:"status"`
	Note      string              `json:"note"`
	Items     []orderItemResponse `json:"items"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

type listOrdersResponse struct {
	Data    []orderResponse `json:"data"`
	Page    int32           `json:"page"`
	PerPage int32           `json:"per_page"`
}

func (h *OrderHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListOrdersInput{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]orderResponse, len(result.Data))
	for i, o := range result.Data {
		resp[i] = toOrderResponse(o)
	}

	writeJSON(w, http.StatusOK, listOrdersResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *OrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	order, err := h.service.Create(r.Context(), application.CreateOrderInput{
		Reference: req.Reference,
		Note:      req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toOrderResponse(order))
}

func (h *OrderHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOrderResponse(order))
}

func (h *OrderHandler) addItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req addOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.AddItem(r.Context(), application.AddOrderItemInput{
		OrderID:    id,
		ProductID:  req.ProductID,
		LocationID: req.LocationID,
		Quantity:   req.Quantity,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toOrderItemResponse(item))
}

func (h *OrderHandler) confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := h.service.Confirm(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOrderResponse(order))
}

func (h *OrderHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	order, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOrderResponse(order))
}

func toOrderResponse(o *domain.Order) orderResponse {
	items := make([]orderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = toOrderItemResponse(item)
	}

	return orderResponse{
		ID:        o.ID,
		Reference: o.Reference,
		Status:    string(o.Status),
		Note:      o.Note,
		Items:     items,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toOrderItemResponse(item *domain.OrderItem) orderItemResponse {
	return orderItemResponse{
		ID:         item.ID,
		OrderID:    item.OrderID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   item.Quantity,
		CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
		case "REFERENCE_ALREADY_EXISTS", "INVALID_REFERENCE", "INVALID_STATUS",
			"INVALID_ORDER", "INVALID_PRODUCT_ID", "INVALID_LOCATION_ID", "INVALID_QUANTITY":
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
