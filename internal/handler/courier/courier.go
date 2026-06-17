package courier

import (
	modelCourier "avito/internal/model/courier"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service service
}

func NewHandler(svc service) *Handler {
	return &Handler{svc}
}

func (h *Handler) GetCourier(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	courier, err := h.service.GetByID(context.Background(), id)
	if err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusOK, modelToResponse(courier))
}

func (h *Handler) CreateCourier(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := &CreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := req.validate(); err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	id, err := h.service.Create(context.Background(), req.toModel())
	if err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id": id,
		"msg": "courier created",
	})
}

func (h *Handler) UpdateCourier(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req := &UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := req.validate(); err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	if err := h.service.Update(context.Background(), req.toModel()); err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.service.GetAll(context.Background())
	if err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	respCouriers := make([]Courier, len(couriers))
	for i, v := range couriers {
		respCouriers[i] = modelToResponse(v)
	}
	writeJSON(w, http.StatusOK, respCouriers)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	if err := h.service.Delete(context.Background(), id); err != nil {
		status, msg := mapErrorToHTTP(err)
		writeError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusNoContent, map[string]any{"msg": "courier deleted"})
}

func getID(r *http.Request) (int, error) {
	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		return -1, modelCourier.ErrInvalidId
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return -1, modelCourier.ErrInvalidId
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
