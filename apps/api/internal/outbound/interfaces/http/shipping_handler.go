package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type ShippingHandler struct {
	service *application.ShippingService
}

func NewShippingHandler(service *application.ShippingService) *ShippingHandler {
	return &ShippingHandler{service: service}
}

func (h *ShippingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/shippings", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/ship", h.ship)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type shipRequest struct {
	TrackingCode string `json:"tracking_code"`
}

type shippingResponse struct {
	ID           string `json:"id"`
	OrderID      string `json:"order_id"`
	PackingID    string `json:"packing_id"`
	Status       string `json:"status"`
	TrackingCode string `json:"tracking_code"`
	Note         string `json:"note"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (h *ShippingHandler) create(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order_id is required"})
		return
	}

	shipping, err := h.service.Create(r.Context(), orderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toShippingResponse(shipping))
}

func (h *ShippingHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	shipping, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toShippingResponse(shipping))
}

func (h *ShippingHandler) ship(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req shipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	shipping, err := h.service.Ship(r.Context(), id, req.TrackingCode)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toShippingResponse(shipping))
}

func (h *ShippingHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	shipping, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toShippingResponse(shipping))
}

func toShippingResponse(s *domain.Shipping) shippingResponse {
	return shippingResponse{
		ID:           s.ID,
		OrderID:      s.OrderID,
		PackingID:    s.PackingID,
		Status:       string(s.Status),
		TrackingCode: s.TrackingCode,
		Note:         s.Note,
		CreatedAt:    s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
