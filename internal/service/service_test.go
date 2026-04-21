package service

import (
	"context"
	"errors"
	"testing"

	"personal-manager/internal/model"
	"personal-manager/internal/store"
)

func TestCreateValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		person  model.Person
		wantErr string
	}{
		{
			name: "missing userid",
			person: model.Person{
				Name:  "Alice",
				Email: "alice@example.com",
				Phone: "13800138000",
			},
			wantErr: "userid is required",
		},
		{
			name: "missing name",
			person: model.Person{
				UserID: "u1",
				Email:  "alice@example.com",
				Phone:  "13800138000",
			},
			wantErr: "name is required",
		},
		{
			name: "missing email",
			person: model.Person{
				UserID: "u1",
				Name:   "Alice",
				Phone:  "13800138000",
			},
			wantErr: "email is required",
		},
		{
			name: "missing phone",
			person: model.Person{
				UserID: "u1",
				Name:   "Alice",
				Email:  "alice@example.com",
			},
			wantErr: "phone is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(&fakeStore{})

			_, err := svc.Create(context.Background(), tt.person)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("Create() error = %v, want validation error", err)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("Create() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateNormalizesAndStoresPerson(t *testing.T) {
	st := &fakeStore{}
	svc := New(st)

	got, err := svc.Create(context.Background(), model.Person{
		UserID: " u1 ",
		Name:   " Alice ",
		Email:  " alice@example.com ",
		Phone:  " 13800138000 ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	want := model.Person{
		UserID: "u1",
		Name:   "Alice",
		Email:  "alice@example.com",
		Phone:  "13800138000",
	}
	if got != want {
		t.Fatalf("Create() = %#v, want %#v", got, want)
	}
	if st.created != want {
		t.Fatalf("stored person = %#v, want %#v", st.created, want)
	}
}

func TestCreateValidatesUserIDFormat(t *testing.T) {
	tests := []struct {
		name   string
		userid string
	}{
		{name: "starts with digit", userid: "1user"},
		{name: "contains hyphen", userid: "u-1"},
		{name: "contains underscore", userid: "u_1"},
		{name: "contains space", userid: "user 1"},
		{name: "contains punctuation", userid: "user!"},
		{name: "contains non ASCII letter", userid: "用户1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(&fakeStore{})

			_, err := svc.Create(context.Background(), model.Person{
				UserID: tt.userid,
				Name:   "Alice",
				Email:  "alice@example.com",
				Phone:  "13800138000",
			})
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("Create() error = %v, want validation error", err)
			}
			if err.Error() != "userid must start with a letter and contain letters and digits only" {
				t.Fatalf("Create() error = %q, want userid validation message", err.Error())
			}
		})
	}
}

func TestCreateValidatesEmailFormat(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{name: "missing at sign", email: "alice.example.com"},
		{name: "missing local part", email: "@example.com"},
		{name: "missing domain", email: "alice@"},
		{name: "contains space", email: "alice smith@example.com"},
		{name: "display name", email: "Alice <alice@example.com>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(&fakeStore{})

			_, err := svc.Create(context.Background(), model.Person{
				UserID: "u1",
				Name:   "Alice",
				Email:  tt.email,
				Phone:  "13800138000",
			})
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("Create() error = %v, want validation error", err)
			}
			if err.Error() != "email must be a valid email address" {
				t.Fatalf("Create() error = %q, want email validation message", err.Error())
			}
		})
	}
}

func TestCreateValidatesMainlandChinaMobilePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
	}{
		{name: "invalid prefix", phone: "12800138000"},
		{name: "too short", phone: "1380013800"},
		{name: "too long", phone: "138001380000"},
		{name: "contains non digit", phone: "1380013800a"},
		{name: "with country code", phone: "+8613800138000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(&fakeStore{})

			_, err := svc.Create(context.Background(), model.Person{
				UserID: "u1",
				Name:   "Alice",
				Email:  "alice@example.com",
				Phone:  tt.phone,
			})
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("Create() error = %v, want validation error", err)
			}
			if err.Error() != "phone must be a valid mainland China mobile number" {
				t.Fatalf("Create() error = %q, want phone validation message", err.Error())
			}
		})
	}
}

func TestCreateDuplicateUserIDReturnsValidationError(t *testing.T) {
	svc := New(&fakeStore{createErr: store.ErrDuplicate})

	_, err := svc.Create(context.Background(), model.Person{
		UserID: "u1",
		Name:   "Alice",
		Email:  "alice@example.com",
		Phone:  "13800138000",
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("Create() error = %v, want validation error", err)
	}
	if err.Error() != "userid already exists" {
		t.Fatalf("Create() error = %q, want duplicate message", err.Error())
	}
}

func TestReadAndDeleteValidateUserID(t *testing.T) {
	svc := New(&fakeStore{})

	if _, err := svc.Read(context.Background(), " "); !errors.Is(err, ErrValidation) {
		t.Fatalf("Read() error = %v, want validation error", err)
	}

	if err := svc.Delete(context.Background(), " "); !errors.Is(err, ErrValidation) {
		t.Fatalf("Delete() error = %v, want validation error", err)
	}
}

func TestReadAndDeleteValidateUserIDFormat(t *testing.T) {
	svc := New(&fakeStore{})

	if _, err := svc.Read(context.Background(), "1user"); !errors.Is(err, ErrValidation) {
		t.Fatalf("Read() error = %v, want validation error", err)
	}

	if err := svc.Delete(context.Background(), "user-1"); !errors.Is(err, ErrValidation) {
		t.Fatalf("Delete() error = %v, want validation error", err)
	}
}

type fakeStore struct {
	created   model.Person
	createErr error
}

func (f *fakeStore) Create(_ context.Context, person model.Person) error {
	f.created = person
	return f.createErr
}

func (f *fakeStore) Get(_ context.Context, userid string) (model.Person, error) {
	return model.Person{UserID: userid}, nil
}

func (f *fakeStore) Update(_ context.Context, person model.Person) error {
	return nil
}

func (f *fakeStore) Delete(_ context.Context, userid string) error {
	return nil
}
