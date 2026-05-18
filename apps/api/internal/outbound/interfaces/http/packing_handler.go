package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type PackingHandler struct {
	service *application.PackingService
}

func NewPackingHandler(service *application.PackingService) *PackingHandler {
	return &PackingHandler{service: service}
}

func (h *PackingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/packings", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type packingResponse struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	PickingID string `json:"picking_id"`
	Status    string `json:"status"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *PackingHandler) create(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order_id is required"})
		return
	}

	packing, err := h.service.Create(r.Context(), orderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPackingResponse(packing))
}

func (h *PackingHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	packing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPackingResponse(packing))
}

func (h *PackingHandler) start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	packing, err := h.service.Start(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPackingResponse(packing))
}

func (h *PackingHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	packing, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPackingResponse(packing))
}

func (h *PackingHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	packing, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPackingResponse(packing))
}

func toPackingResponse(p *domain.Packing) packingResponse {
	return packingResponse{
		ID:        p.ID,
		OrderID:   p.OrderID,
		PickingID: p.PickingID,
		Status:    string(p.Status),
		Note:      p.Note,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
