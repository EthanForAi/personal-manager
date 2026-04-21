package service

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"personal-manager/internal/model"
	"personal-manager/internal/store"
)

var ErrValidation = errors.New("validation failed")

var (
	userIDPattern              = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*$`)
	mainlandChinaMobilePattern = regexp.MustCompile(`^1[3-9][0-9]{9}$`)
)

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
	if err := validateUserID(userid); err != nil {
		return model.Person{}, err
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
	if err := validateUserID(userid); err != nil {
		return err
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
	if err := validateUserID(person.UserID); err != nil {
		return err
	}

	switch {
	case person.Name == "":
		return validationError("name is required")
	case person.Email == "":
		return validationError("email is required")
	case !validEmail(person.Email):
		return validationError("email must be a valid email address")
	case person.Phone == "":
		return validationError("phone is required")
	case !mainlandChinaMobilePattern.MatchString(person.Phone):
		return validationError("phone must be a valid mainland China mobile number")
	default:
		return nil
	}
}

func validateUserID(userid string) error {
	switch {
	case userid == "":
		return validationError("userid is required")
	case !userIDPattern.MatchString(userid):
		return validationError("userid must start with a letter and contain letters and digits only")
	default:
		return nil
	}
}

func validEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func validationError(message string) error {
	return ValidationError{Message: message}
}
