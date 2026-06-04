package store

import (
	"context"
	"database/sql"
	"errors"

	"personal-manager/internal/model"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicate      = errors.New("userid already exists")
	ErrDuplicateEmail = errors.New("email already exists")
)

type Store struct {
	db *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	st := &Store{db: db}
	if err := st.init(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return st, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS personal_info (
	userid TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL
)`)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
CREATE UNIQUE INDEX IF NOT EXISTS idx_personal_info_email
ON personal_info(email)`)
	return err
}

func (s *Store) Create(ctx context.Context, person model.Person) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO personal_info (userid, name, email, phone)
VALUES (?, ?, ?, ?)`,
		person.UserID, person.Name, person.Email, person.Phone,
	)
	if err != nil {
		return s.duplicateError(ctx, err, person)
	}

	return nil
}

func (s *Store) duplicateError(ctx context.Context, err error, person model.Person) error {
	if !isConstraintError(err) {
		return err
	}

	exists, existsErr := s.Exists(ctx, person.UserID)
	if existsErr != nil {
		return existsErr
	}
	if exists {
		return ErrDuplicate
	}

	exists, existsErr = s.emailExists(ctx, person.Email)
	if existsErr != nil {
		return existsErr
	}
	if exists {
		return ErrDuplicateEmail
	}

	return err
}

func isConstraintError(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code()&0xff == sqlite3.SQLITE_CONSTRAINT
}

func (s *Store) emailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
SELECT EXISTS(
	SELECT 1
	FROM personal_info
	WHERE email = ?
)`, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *Store) emailExistsForOtherUserID(ctx context.Context, email, userid string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
SELECT EXISTS(
	SELECT 1
	FROM personal_info
	WHERE email = ? AND userid <> ?
)`, email, userid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *Store) Get(ctx context.Context, userid string) (model.Person, error) {
	var person model.Person
	err := s.db.QueryRowContext(ctx, `
SELECT userid, name, email, phone
FROM personal_info
WHERE userid = ?`, userid).Scan(&person.UserID, &person.Name, &person.Email, &person.Phone)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Person{}, ErrNotFound
	}
	if err != nil {
		return model.Person{}, err
	}

	return person, nil
}

func (s *Store) Update(ctx context.Context, person model.Person) error {
	result, err := s.db.ExecContext(ctx, `
UPDATE personal_info
SET name = ?, email = ?, phone = ?
WHERE userid = ?`,
		person.Name, person.Email, person.Phone, person.UserID,
	)
	if err != nil {
		if isConstraintError(err) {
			exists, existsErr := s.emailExistsForOtherUserID(ctx, person.Email, person.UserID)
			if existsErr != nil {
				return existsErr
			}
			if exists {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, userid string) error {
	result, err := s.db.ExecContext(ctx, `
DELETE FROM personal_info
WHERE userid = ?`, userid)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Store) Exists(ctx context.Context, userid string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `
SELECT EXISTS(
	SELECT 1
	FROM personal_info
	WHERE userid = ?
)`, userid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
