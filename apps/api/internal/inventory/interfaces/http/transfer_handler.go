package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type TransferHandler struct {
	service *application.TransferService
}

func NewTransferHandler(service *application.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

func (h *TransferHandler) RegisterRoutes(r chi.Router) {
	r.Route("/inventory/transfers", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type createTransferRequest struct {
	ProductID      string  `json:"product_id"`
	FromLocationID string  `json:"from_location_id"`
	ToLocationID   string  `json:"to_location_id"`
	Quantity       float64 `json:"quantity"`
	Note           string  `json:"note"`
}

type transferResponse struct {
	ID             string  `json:"id"`
	ProductID      string  `json:"product_id"`
	FromLocationID string  `json:"from_location_id"`
	ToLocationID   string  `json:"to_location_id"`
	Quantity       float64 `json:"quantity"`
	Status         string  `json:"status"`
	Note           string  `json:"note"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type listTransfersResponse struct {
	Data    []transferResponse `json:"data"`
	Page    int32              `json:"page"`
	PerPage int32              `json:"per_page"`
}

func (h *TransferHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListTransfersInput{
		ProductID:      r.URL.Query().Get("product_id"),
		FromLocationID: r.URL.Query().Get("from_location_id"),
		ToLocationID:   r.URL.Query().Get("to_location_id"),
		Page:           page,
		PerPage:        perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]transferResponse, len(result.Data))
	for i, t := range result.Data {
		resp[i] = toTransferResponse(t)
	}

	writeJSON(w, http.StatusOK, listTransfersResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *TransferHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	transfer, err := h.service.Create(r.Context(), application.CreateTransferInput{
		ProductID:      req.ProductID,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Quantity:       req.Quantity,
		Note:           req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTransferResponse(transfer))
}

func (h *TransferHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	transfer, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTransferResponse(transfer))
}

func (h *TransferHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	transfer, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTransferResponse(transfer))
}

func (h *TransferHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	transfer, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTransferResponse(transfer))
}

func toTransferResponse(t *domain.StockTransfer) transferResponse {
	return transferResponse{
		ID:             t.ID,
		ProductID:      t.ProductID,
		FromLocationID: t.FromLocationID,
		ToLocationID:   t.ToLocationID,
		Quantity:       t.Quantity,
		Status:         string(t.Status),
		Note:           t.Note,
		CreatedAt:      t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
