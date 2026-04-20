package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"personal-manager/internal/model"
	"personal-manager/internal/service"
	"personal-manager/internal/store"
)

type Service interface {
	Create(ctx context.Context, person model.Person) (model.Person, error)
	Read(ctx context.Context, userid string) (model.Person, error)
	Update(ctx context.Context, person model.Person) (model.Person, error)
	Delete(ctx context.Context, userid string) error
}

type Handler struct {
	service Service
}

type idRequest struct {
	UserID string `json:"userid"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.handleCreate)
	mux.HandleFunc("/read", h.handleRead)
	mux.HandleFunc("/update", h.handleUpdate)
	mux.HandleFunc("/delete", h.handleDelete)
	return mux
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	var req model.Person
	if !decodeJSON(w, r, &req) {
		return
	}

	person, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *Handler) handleRead(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	var req idRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	person, err := h.service.Read(r.Context(), req.UserID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	var req model.Person
	if !decodeJSON(w, r, &req) {
		return
	}

	person, err := h.service.Update(r.Context(), req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, person)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	var req idRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.service.Delete(r.Context(), req.UserID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, model.DeleteResponse{Deleted: true})
}

func requirePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		return true
	}
	writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	return false
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON"})
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrValidation):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, store.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "record not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
