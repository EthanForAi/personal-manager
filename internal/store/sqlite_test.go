package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"personal-manager/internal/model"
)

func TestStoreCreateGetUpdateDelete(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	person := model.Person{
		UserID: "u-1",
		Name:   "Alice",
		Email:  "alice@example.com",
		Phone:  "13800138000",
	}
	if err := st.Create(ctx, person); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := st.Get(ctx, person.UserID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != person {
		t.Fatalf("Get() = %#v, want %#v", got, person)
	}

	if err := st.Create(ctx, person); !errors.Is(err, ErrDuplicate) {
		t.Fatalf("Create() duplicate error = %v, want %v", err, ErrDuplicate)
	}

	updated := model.Person{
		UserID: "u-1",
		Name:   "Alice Smith",
		Email:  "alice.smith@example.com",
		Phone:  "13900139000",
	}
	if err := st.Update(ctx, updated); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err = st.Get(ctx, updated.UserID)
	if err != nil {
		t.Fatalf("Get() after update error = %v", err)
	}
	if got != updated {
		t.Fatalf("Get() after update = %#v, want %#v", got, updated)
	}

	if err := st.Delete(ctx, updated.UserID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	if _, err := st.Get(ctx, updated.UserID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() deleted error = %v, want %v", err, ErrNotFound)
	}
}

func TestStoreMissingRecords(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.Get(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() missing error = %v, want %v", err, ErrNotFound)
	}

	err := st.Update(ctx, model.Person{
		UserID: "missing",
		Name:   "Missing",
		Email:  "missing@example.com",
		Phone:  "13800138000",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update() missing error = %v, want %v", err, ErrNotFound)
	}

	if err := st.Delete(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete() missing error = %v, want %v", err, ErrNotFound)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()

	st, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	return st
}
