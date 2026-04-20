package service

import (
	"context"
	"errors"
	"strings"

	"personal-manager/internal/model"
	"personal-manager/internal/store"
)

var ErrValidation = errors.New("validation failed")

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) Is(target error) bool {
	return target == ErrValidation
}

type PersonStore interface {
	Create(context.Context, model.Person) error
	Get(context.Context, string) (model.Person, error)
	Update(context.Context, model.Person) error
	Delete(context.Context, string) error
}

type Service struct {
	store PersonStore
}

func New(store PersonStore) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, person model.Person) (model.Person, error) {
	person = normalize(person)
	if err := validatePerson(person); err != nil {
		return model.Person{}, err
	}

	if err := s.store.Create(ctx, person); err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			return model.Person{}, validationError("userid already exists")
		}
		return model.Person{}, err
	}

	return person, nil
}

func (s *Service) Read(ctx context.Context, userid string) (model.Person, error) {
	userid = strings.TrimSpace(userid)
	if userid == "" {
		return model.Person{}, validationError("userid is required")
	}

	return s.store.Get(ctx, userid)
}

func (s *Service) Update(ctx context.Context, person model.Person) (model.Person, error) {
	person = normalize(person)
	if err := validatePerson(person); err != nil {
		return model.Person{}, err
	}

	if err := s.store.Update(ctx, person); err != nil {
		return model.Person{}, err
	}

	return person, nil
}

func (s *Service) Delete(ctx context.Context, userid string) error {
	userid = strings.TrimSpace(userid)
	if userid == "" {
		return validationError("userid is required")
	}

	return s.store.Delete(ctx, userid)
}

func normalize(person model.Person) model.Person {
	person.UserID = strings.TrimSpace(person.UserID)
	person.Name = strings.TrimSpace(person.Name)
	person.Email = strings.TrimSpace(person.Email)
	person.Phone = strings.TrimSpace(person.Phone)
	return person
}

func validatePerson(person model.Person) error {
	switch {
	case person.UserID == "":
		return validationError("userid is required")
	case person.Name == "":
		return validationError("name is required")
	case person.Email == "":
		return validationError("email is required")
	case person.Phone == "":
		return validationError("phone is required")
	default:
		return nil
	}
}

func validationError(message string) error {
	return ValidationError{Message: message}
}
