package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	dderr "github.com/aleodoni/go-ddd/errors"
	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type PurchaseOrderHandler struct {
	service *application.PurchaseOrderService
}

func NewPurchaseOrderHandler(service *application.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{service: service}
}

func (h *PurchaseOrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/purchase-orders", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/items", h.addItem)
		r.Post("/{id}/confirm", h.confirm)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type createPurchaseOrderRequest struct {
	Reference string `json:"reference"`
	Supplier  string `json:"supplier"`
	Note      string `json:"note"`
}

type addPurchaseOrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type purchaseOrderItemResponse struct {
	ID              string  `json:"id"`
	PurchaseOrderID string  `json:"purchase_order_id"`
	ProductID       string  `json:"product_id"`
	Quantity        float64 `json:"quantity"`
	CreatedAt       string  `json:"created_at"`
}

type purchaseOrderResponse struct {
	ID        string                      `json:"id"`
	Reference string                      `json:"reference"`
	Supplier  string                      `json:"supplier"`
	Status    string                      `json:"status"`
	Note      string                      `json:"note"`
	Items     []purchaseOrderItemResponse `json:"items"`
	CreatedAt string                      `json:"created_at"`
	UpdatedAt string                      `json:"updated_at"`
}

type listPurchaseOrdersResponse struct {
	Data    []purchaseOrderResponse `json:"data"`
	Page    int32                   `json:"page"`
	PerPage int32                   `json:"per_page"`
}

func (h *PurchaseOrderHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListPurchaseOrdersInput{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]purchaseOrderResponse, len(result.Data))
	for i, po := range result.Data {
		resp[i] = toPurchaseOrderResponse(po)
	}

	writeJSON(w, http.StatusOK, listPurchaseOrdersResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *PurchaseOrderHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createPurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	po, err := h.service.Create(r.Context(), application.CreatePurchaseOrderInput{
		Reference: req.Reference,
		Supplier:  req.Supplier,
		Note:      req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPurchaseOrderResponse(po))
}

func (h *PurchaseOrderHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	po, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPurchaseOrderResponse(po))
}

func (h *PurchaseOrderHandler) addItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req addPurchaseOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.AddItem(r.Context(), application.AddPurchaseOrderItemInput{
		PurchaseOrderID: id,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPurchaseOrderItemResponse(item))
}

func (h *PurchaseOrderHandler) confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	po, err := h.service.Confirm(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPurchaseOrderResponse(po))
}

func (h *PurchaseOrderHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	po, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPurchaseOrderResponse(po))
}

func toPurchaseOrderResponse(po *domain.PurchaseOrder) purchaseOrderResponse {
	items := make([]purchaseOrderItemResponse, len(po.Items))
	for i, item := range po.Items {
		items[i] = toPurchaseOrderItemResponse(item)
	}

	return purchaseOrderResponse{
		ID:        po.ID,
		Reference: po.Reference,
		Supplier:  po.Supplier,
		Status:    string(po.Status),
		Note:      po.Note,
		Items:     items,
		CreatedAt: po.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: po.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPurchaseOrderItemResponse(item *domain.PurchaseOrderItem) purchaseOrderItemResponse {
	return purchaseOrderItemResponse{
		ID:              item.ID,
		PurchaseOrderID: item.PurchaseOrderID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		CreatedAt:       item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
			"INVALID_PURCHASE_ORDER", "INVALID_PRODUCT_ID", "INVALID_QUANTITY",
			"INVALID_LOCATION_ID", "INCOMPLETE_PUT_AWAY", "INVALID_PUT_AWAY_ID",
			"INVALID_RECEIPT_ID", "INVALID_PURCHASE_ORDER_ID":
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
