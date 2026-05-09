package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/application"
	"github.com/aleooni-fullstack/ale-wms/apps/api/internal/location/domain"
)

type ZoneHandler struct {
	service *application.ZoneService
}

func NewZoneHandler(service *application.ZoneService) *ZoneHandler {
	return &ZoneHandler{service: service}
}

func (h *ZoneHandler) RegisterRoutes(r chi.Router) {
	r.Route("/warehouses/{warehouseID}/zones", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.deactivate)
	})
}

type createZoneRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type updateZoneRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type zoneResponse struct {
	ID          string `json:"id"`
	WarehouseID string `json:"warehouse_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
}

type listZonesResponse struct {
	Data    []zoneResponse `json:"data"`
	Page    int32          `json:"page"`
	PerPage int32          `json:"per_page"`
}

func (h *ZoneHandler) list(w http.ResponseWriter, r *http.Request) {
	warehouseID := chi.URLParam(r, "warehouseID")
	page, perPage := parsePagination(r)

	result, err := h.service.List(r.Context(), application.ListZonesInput{
		WarehouseID: warehouseID,
		Page:        page,
		PerPage:     perPage,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]zoneResponse, len(result.Data))
	for i, z := range result.Data {
		resp[i] = toZoneResponse(z)
	}

	writeJSON(w, http.StatusOK, listZonesResponse{
		Data:    resp,
		Page:    result.Page,
		PerPage: result.PerPage,
	})
}

func (h *ZoneHandler) create(w http.ResponseWriter, r *http.Request) {
	warehouseID := chi.URLParam(r, "warehouseID")

	var req createZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	zone, err := h.service.Create(r.Context(), application.CreateZoneInput{
		WarehouseID: warehouseID,
		Code:        req.Code,
		Name:        req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toZoneResponse(zone))
}

func (h *ZoneHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	zone, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toZoneResponse(zone))
}

func (h *ZoneHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	zone, err := h.service.Update(r.Context(), id, application.UpdateZoneInput{
		Code: req.Code,
		Name: req.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toZoneResponse(zone))
}

func (h *ZoneHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Deactivate(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toZoneResponse(z *domain.Zone) zoneResponse {
	return zoneResponse{
		ID:          z.ID,
		WarehouseID: z.WarehouseID,
		Code:        z.Code,
		Name:        z.Name,
		Active:      z.Active,
	}
}
