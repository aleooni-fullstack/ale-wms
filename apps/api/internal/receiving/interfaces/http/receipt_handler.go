package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type ReceiptHandler struct {
	service *application.ReceiptService
}

func NewReceiptHandler(service *application.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

func (h *ReceiptHandler) RegisterRoutes(r chi.Router) {
	r.Route("/receipts", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
		r.Post("/{id}/items/{productID}/receive", h.receiveItem)
	})
}

type receiveItemRequest struct {
	ReceivedQuantity float64 `json:"received_quantity"`
}

type receiptItemResponse struct {
	ID               string   `json:"id"`
	ReceiptID        string   `json:"receipt_id"`
	ProductID        string   `json:"product_id"`
	ExpectedQuantity float64  `json:"expected_quantity"`
	ReceivedQuantity *float64 `json:"received_quantity"`
	Difference       *float64 `json:"difference"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type receiptResponse struct {
	ID              string                `json:"id"`
	PurchaseOrderID string                `json:"purchase_order_id"`
	Status          string                `json:"status"`
	Note            string                `json:"note"`
	Items           []receiptItemResponse `json:"items"`
	CreatedAt       string                `json:"created_at"`
	UpdatedAt       string                `json:"updated_at"`
}

func (h *ReceiptHandler) create(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.URL.Query().Get("purchase_order_id")
	if purchaseOrderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "purchase_order_id is required"})
		return
	}

	receipt, err := h.service.Create(r.Context(), purchaseOrderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toReceiptResponse(receipt))
}

func (h *ReceiptHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	receipt, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReceiptResponse(receipt))
}

func (h *ReceiptHandler) start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	receipt, err := h.service.Start(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReceiptResponse(receipt))
}

func (h *ReceiptHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	receipt, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReceiptResponse(receipt))
}

func (h *ReceiptHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	receipt, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReceiptResponse(receipt))
}

func (h *ReceiptHandler) receiveItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID := chi.URLParam(r, "productID")

	var req receiveItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.ReceiveItem(r.Context(), application.ReceiveItemInput{
		ReceiptID:        id,
		ProductID:        productID,
		ReceivedQuantity: req.ReceivedQuantity,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReceiptItemResponse(item))
}

func toReceiptResponse(r *domain.Receipt) receiptResponse {
	items := make([]receiptItemResponse, len(r.Items))
	for i, item := range r.Items {
		items[i] = toReceiptItemResponse(item)
	}

	return receiptResponse{
		ID:              r.ID,
		PurchaseOrderID: r.PurchaseOrderID,
		Status:          string(r.Status),
		Note:            r.Note,
		Items:           items,
		CreatedAt:       r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toReceiptItemResponse(item *domain.ReceiptItem) receiptItemResponse {
	return receiptItemResponse{
		ID:               item.ID,
		ReceiptID:        item.ReceiptID,
		ProductID:        item.ProductID,
		ExpectedQuantity: item.ExpectedQuantity,
		ReceivedQuantity: item.ReceivedQuantity,
		Difference:       item.Difference(),
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
