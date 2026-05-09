package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type LocationHandler struct {
	service *application.LocationService
}

func NewLocationHandler(service *application.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/zones/{zoneID}/locations", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.deactivate)
	})
}

type createLocationRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type updateLocationRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type locationResponse struct {
	ID     string `json:"id"`
	ZoneID string `json:"zone_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type listLocationsResponse struct {
	Data    []locationResponse `json:"data"`
	Page    int32              `json:"page"`
	PerPage int32              `json:"per_page"`
}

func (h *LocationHandler) list(w http.ResponseWriter, r *http.Request) {
	zoneID := chi.URLParam(r, "zoneID")
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListLocationsInput{
		ZoneID:  zoneID,
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]locationResponse, len(result.Data))
	for i, l := range result.Data {
		resp[i] = toLocationResponse(l)
	}

	writeJSON(w, http.StatusOK, listLocationsResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *LocationHandler) create(w http.ResponseWriter, r *http.Request) {
	zoneID := chi.URLParam(r, "zoneID")

	var req createLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	location, err := h.service.Create(r.Context(), application.CreateLocationInput{
		ZoneID: zoneID,
		Code:   req.Code,
		Name:   req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toLocationResponse(location))
}

func (h *LocationHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	location, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toLocationResponse(location))
}

func (h *LocationHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	location, err := h.service.Update(r.Context(), id, application.UpdateLocationInput{
		Code: req.Code,
		Name: req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toLocationResponse(location))
}

func (h *LocationHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Deactivate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toLocationResponse(l *domain.Location) locationResponse {
	return locationResponse{
		ID:     l.ID,
		ZoneID: l.ZoneID,
		Code:   l.Code,
		Name:   l.Name,
		Active: l.Active,
	}
}
