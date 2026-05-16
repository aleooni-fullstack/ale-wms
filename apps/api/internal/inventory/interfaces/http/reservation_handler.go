package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/inventory/domain"
)

type ReservationHandler struct {
	service *application.ReservationService
}

func NewReservationHandler(service *application.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/inventory/reservations", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Post("/{id}/confirm", h.confirm)
		r.Post("/{id}/fulfill", h.fulfill)
		r.Post("/{id}/release", h.release)
		r.Post("/{id}/cancel", h.cancel)
	})
}

type createReservationRequest struct {
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	Reference  string  `json:"reference"`
	Note       string  `json:"note"`
}

type reservationResponse struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	LocationID string  `json:"location_id"`
	Quantity   float64 `json:"quantity"`
	Status     string  `json:"status"`
	Reference  string  `json:"reference"`
	Note       string  `json:"note"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type listReservationsResponse struct {
	Data    []reservationResponse `json:"data"`
	Page    int32                 `json:"page"`
	PerPage int32                 `json:"per_page"`
}

func (h *ReservationHandler) list(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListReservationsInput{
		ProductID:  r.URL.Query().Get("product_id"),
		LocationID: r.URL.Query().Get("location_id"),
		Reference:  r.URL.Query().Get("reference"),
		Page:       page,
		PerPage:    perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]reservationResponse, len(result.Data))
	for i, res := range result.Data {
		resp[i] = toReservationResponse(res)
	}

	writeJSON(w, http.StatusOK, listReservationsResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	reservation, err := h.service.Create(r.Context(), application.CreateReservationInput{
		ProductID:  req.ProductID,
		LocationID: req.LocationID,
		Quantity:   req.Quantity,
		Reference:  req.Reference,
		Note:       req.Note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toReservationResponse(reservation))
}

func (h *ReservationHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	reservation, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (h *ReservationHandler) confirm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	reservation, err := h.service.Confirm(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (h *ReservationHandler) fulfill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	reservation, err := h.service.Fulfill(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (h *ReservationHandler) release(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	reservation, err := h.service.Release(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func (h *ReservationHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	reservation, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toReservationResponse(reservation))
}

func toReservationResponse(res *domain.StockReservation) reservationResponse {
	return reservationResponse{
		ID:         res.ID,
		ProductID:  res.ProductID,
		LocationID: res.LocationID,
		Quantity:   res.Quantity,
		Status:     string(res.Status),
		Reference:  res.Reference,
		Note:       res.Note,
		CreatedAt:  res.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  res.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
