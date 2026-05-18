package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/outbound/domain"
)

type PickingHandler struct {
	service *application.PickingService
}

func NewPickingHandler(service *application.PickingService) *PickingHandler {
	return &PickingHandler{service: service}
}

func (h *PickingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/pickings", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
		r.Post("/{id}/items/{productID}/pick", h.pickItem)
	})
}

type pickingItemResponse struct {
	ID         string  `json:"id"`
	PickingID  string  `json:"picking_id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	Picked     bool    `json:"picked"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type pickingResponse struct {
	ID        string                `json:"id"`
	OrderID   string                `json:"order_id"`
	Status    string                `json:"status"`
	Note      string                `json:"note"`
	Items     []pickingItemResponse `json:"items"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

func (h *PickingHandler) create(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "order_id is required"})
		return
	}

	picking, err := h.service.Create(r.Context(), orderID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPickingResponse(picking))
}

func (h *PickingHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	picking, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPickingResponse(picking))
}

func (h *PickingHandler) start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	picking, err := h.service.Start(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPickingResponse(picking))
}

func (h *PickingHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	picking, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPickingResponse(picking))
}

func (h *PickingHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	picking, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPickingResponse(picking))
}

func (h *PickingHandler) pickItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID := chi.URLParam(r, "productID")

	item, err := h.service.PickItem(r.Context(), id, productID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPickingItemResponse(item))
}

func toPickingResponse(p *domain.Picking) pickingResponse {
	items := make([]pickingItemResponse, len(p.Items))
	for i, item := range p.Items {
		items[i] = toPickingItemResponse(item)
	}

	return pickingResponse{
		ID:        p.ID,
		OrderID:   p.OrderID,
		Status:    string(p.Status),
		Note:      p.Note,
		Items:     items,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPickingItemResponse(item *domain.PickingItem) pickingItemResponse {
	return pickingItemResponse{
		ID:         item.ID,
		PickingID:  item.PickingID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   item.Quantity,
		Picked:     item.Picked,
		CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
