package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"personal-manager/internal/model"
	"personal-manager/internal/service"
	"personal-manager/internal/store"
)

func TestHandlerCRUD(t *testing.T) {
	router := newTestRouter(t)

	createBody := `{"userid":"u1","name":"Alice","email":"alice@example.com","phone":"13800138000"}`
	rec := postJSON(router, "/create", createBody)
	assertStatus(t, rec, http.StatusOK)
	assertPerson(t, rec, model.Person{
		UserID: "u1",
		Name:   "Alice",
		Email:  "alice@example.com",
		Phone:  "13800138000",
	})

	rec = postJSON(router, "/read", `{"userid":"u1"}`)
	assertStatus(t, rec, http.StatusOK)
	assertPerson(t, rec, model.Person{
		UserID: "u1",
		Name:   "Alice",
		Email:  "alice@example.com",
		Phone:  "13800138000",
	})

	updateBody := `{"userid":"u1","name":"Alice Smith","email":"alice.smith@example.com","phone":"13900139000"}`
	rec = postJSON(router, "/update", updateBody)
	assertStatus(t, rec, http.StatusOK)
	assertPerson(t, rec, model.Person{
		UserID: "u1",
		Name:   "Alice Smith",
		Email:  "alice.smith@example.com",
		Phone:  "13900139000",
	})

	rec = postJSON(router, "/delete", `{"userid":"u1"}`)
	assertStatus(t, rec, http.StatusOK)
	var deleted model.DeleteResponse
	decodeBody(t, rec, &deleted)
	if !deleted.Deleted {
		t.Fatalf("deleted = false, want true")
	}

	rec = postJSON(router, "/read", `{"userid":"u1"}`)
	assertStatus(t, rec, http.StatusNotFound)
	assertError(t, rec, "record not found")
}

func TestHandlerErrors(t *testing.T) {
	router := newTestRouter(t)

	tests := []struct {
		name      string
		method    string
		path      string
		body      string
		wantCode  int
		wantError string
	}{
		{
			name:      "invalid JSON",
			method:    http.MethodPost,
			path:      "/create",
			body:      `{`,
			wantCode:  http.StatusBadRequest,
			wantError: "invalid JSON",
		},
		{
			name:      "validation failure",
			method:    http.MethodPost,
			path:      "/create",
			body:      `{"name":"Alice","email":"alice@example.com","phone":"13800138000"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "userid is required",
		},
		{
			name:      "invalid userid",
			method:    http.MethodPost,
			path:      "/create",
			body:      `{"userid":"1user","name":"Alice","email":"alice@example.com","phone":"13800138000"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "userid must start with a letter and contain letters and digits only",
		},
		{
			name:      "invalid email",
			method:    http.MethodPost,
			path:      "/create",
			body:      `{"userid":"u1","name":"Alice","email":"alice.example.com","phone":"13800138000"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "email must be a valid email address",
		},
		{
			name:      "invalid phone",
			method:    http.MethodPost,
			path:      "/create",
			body:      `{"userid":"u1","name":"Alice","email":"alice@example.com","phone":"12800138000"}`,
			wantCode:  http.StatusBadRequest,
			wantError: "phone must be a valid mainland China mobile number",
		},
		{
			name:      "missing record",
			method:    http.MethodPost,
			path:      "/read",
			body:      `{"userid":"missing"}`,
			wantCode:  http.StatusNotFound,
			wantError: "record not found",
		},
		{
			name:      "non POST method",
			method:    http.MethodGet,
			path:      "/create",
			wantCode:  http.StatusMethodNotAllowed,
			wantError: "method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assertStatus(t, rec, tt.wantCode)
			assertError(t, rec, tt.wantError)
		})
	}
}

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()

	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("store.Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("store.Close() error = %v", err)
		}
	})

	return New(service.New(st)).Routes()
}

func postJSON(handler http.Handler, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	return serve(handler, req)
}

func serve(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, want, rec.Body.String())
	}
}

func assertPerson(t *testing.T, rec *httptest.ResponseRecorder, want model.Person) {
	t.Helper()

	var got model.Person
	decodeBody(t, rec, &got)
	if got != want {
		t.Fatalf("person = %#v, want %#v", got, want)
	}
}

func assertError(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()

	var got struct {
		Error string `json:"error"`
	}
	decodeBody(t, rec, &got)
	if got.Error != want {
		t.Fatalf("error = %q, want %q", got.Error, want)
	}
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	if err := json.NewDecoder(rec.Body).Decode(dst); err != nil {
		t.Fatalf("decode body error = %v; body = %s", err, rec.Body.String())
	}
}
