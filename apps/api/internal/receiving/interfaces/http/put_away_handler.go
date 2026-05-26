package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/receiving/domain"
)

type PutAwayHandler struct {
	service *application.PutAwayService
}

func NewPutAwayHandler(service *application.PutAwayService) *PutAwayHandler {
	return &PutAwayHandler{service: service}
}

func (h *PutAwayHandler) RegisterRoutes(r chi.Router) {
	r.Route("/put-aways", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/start", h.start)
		r.Post("/{id}/complete", h.complete)
		r.Post("/{id}/cancel", h.cancel)
		r.Post("/{id}/items", h.addItem)
		r.Post("/{id}/items/{productID}/store", h.storeItem)
	})
}

type addPutAwayItemRequest struct {
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
}

type putAwayItemResponse struct {
	ID         string  `json:"id"`
	PutAwayID  string  `json:"put_away_id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	PutAway    bool    `json:"put_away"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type putAwayResponse struct {
	ID        string                `json:"id"`
	ReceiptID string                `json:"receipt_id"`
	Status    string                `json:"status"`
	Note      string                `json:"note"`
	Items     []putAwayItemResponse `json:"items"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

func (h *PutAwayHandler) create(w http.ResponseWriter, r *http.Request) {
	receiptID := r.URL.Query().Get("receipt_id")
	if receiptID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "receipt_id is required"})
		return
	}

	putAway, err := h.service.Create(r.Context(), receiptID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPutAwayResponse(putAway))
}

func (h *PutAwayHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	putAway, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPutAwayResponse(putAway))
}

func (h *PutAwayHandler) start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	putAway, err := h.service.Start(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPutAwayResponse(putAway))
}

func (h *PutAwayHandler) complete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	putAway, err := h.service.Complete(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPutAwayResponse(putAway))
}

func (h *PutAwayHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	putAway, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPutAwayResponse(putAway))
}

func (h *PutAwayHandler) addItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req addPutAwayItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	item, err := h.service.AddItem(r.Context(), application.AddPutAwayItemInput{
		PutAwayID:  id,
		ProductID:  req.ProductID,
		LocationID: req.LocationID,
		Quantity:   req.Quantity,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toPutAwayItemResponse(item))
}

func (h *PutAwayHandler) storeItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID := chi.URLParam(r, "productID")

	item, err := h.service.StoreItem(r.Context(), application.StoreItemInput{
		PutAwayID: id,
		ProductID: productID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toPutAwayItemResponse(item))
}

func toPutAwayResponse(p *domain.PutAway) putAwayResponse {
	items := make([]putAwayItemResponse, len(p.Items))
	for i, item := range p.Items {
		items[i] = toPutAwayItemResponse(item)
	}

	return putAwayResponse{
		ID:        p.ID,
		ReceiptID: p.ReceiptID,
		Status:    string(p.Status),
		Note:      p.Note,
		Items:     items,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPutAwayItemResponse(item *domain.PutAwayItem) putAwayItemResponse {
	return putAwayItemResponse{
		ID:         item.ID,
		PutAwayID:  item.PutAwayID,
		ProductID:  item.ProductID,
		LocationID: item.LocationID,
		Quantity:   item.Quantity,
		PutAway:    item.PutAway,
		CreatedAt:  item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
