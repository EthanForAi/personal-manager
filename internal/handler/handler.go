package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"personal-manager/internal/model"
	"personal-manager/internal/service"
	"personal-manager/internal/store"
)

type Service interface {
	Create(ctx context.Context, person model.Person) (model.Person, error)
	Read(ctx context.Context, userid string) (model.Person, error)
	Update(ctx context.Context, person model.Person) (model.Person, error)
	Delete(ctx context.Context, userid string) error
	Check(ctx context.Context, userid string) (bool, error)
}

type Handler struct {
	service Service
	logger  *log.Logger
	now     func() time.Time
}

type idRequest struct {
	UserID string `json:"userid"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func New(service Service) *Handler {
	return NewWithLogger(service, log.Default())
}

func NewWithLogger(service Service, logger *log.Logger) *Handler {
	return newWithLoggerAndClock(service, logger, time.Now)
}

func newWithLoggerAndClock(service Service, logger *log.Logger, now func() time.Time) *Handler {
	if logger == nil {
		logger = log.Default()
	}
	if now == nil {
		now = time.Now
	}
	return &Handler{service: service, logger: logger, now: now}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", h.handleCreate)
	mux.HandleFunc("/read", h.handleRead)
	mux.HandleFunc("/update", h.handleUpdate)
	mux.HandleFunc("/delete", h.handleDelete)
	mux.HandleFunc("/check", h.handleCheck)
	return h.logRequests(mux)
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
		h.logOperation("create", req.UserID, err)
		writeError(w, err)
		return
	}

	h.logOperation("create", person.UserID, nil)
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
		h.logOperation("read", req.UserID, err)
		writeError(w, err)
		return
	}

	h.logOperation("read", person.UserID, nil)
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
		h.logOperation("update", req.UserID, err)
		writeError(w, err)
		return
	}

	h.logOperation("update", person.UserID, nil)
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
		h.logOperation("delete", req.UserID, err)
		writeError(w, err)
		return
	}

	h.logOperation("delete", req.UserID, nil)
	writeJSON(w, http.StatusOK, model.DeleteResponse{Deleted: true})
}

func (h *Handler) handleCheck(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}

	var req idRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	exists, err := h.service.Check(r.Context(), req.UserID)
	if err != nil {
		h.logOperation("check", req.UserID, err)
		writeError(w, err)
		return
	}

	h.logOperation("check", req.UserID, nil)
	writeJSON(w, http.StatusOK, model.CheckResponse{Exists: exists})
}

func (h *Handler) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := h.now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		h.logf(
			"request completed method=%s path=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			rec.status,
			h.now().Sub(start),
		)
	})
}

func (h *Handler) logOperation(operation string, userid string, err error) {
	if err != nil {
		h.logf("operation=%s userid=%q result=error error=%q", operation, userid, err.Error())
		return
	}
	h.logf("operation=%s userid=%q result=success", operation, userid)
}

func (h *Handler) logf(format string, args ...any) {
	logArgs := make([]any, 0, len(args)+1)
	logArgs = append(logArgs, h.now().UTC().Format(time.RFC3339Nano))
	logArgs = append(logArgs, args...)
	h.logger.Printf("timestamp=%s "+format, logArgs...)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(body)
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
