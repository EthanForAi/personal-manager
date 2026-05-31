package store

import (
	"context"
	"database/sql"
	"errors"

	"personal-manager/internal/model"

	_ "modernc.org/sqlite"
)

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("userid already exists")
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
	email TEXT NOT NULL,
	phone TEXT NOT NULL,
	introduction TEXT NOT NULL DEFAULT ''
)`)
	if err != nil {
		return err
	}

	return s.ensureIntroductionColumn(ctx)
}

func (s *Store) ensureIntroductionColumn(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `PRAGMA table_info(personal_info)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == "introduction" {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `ALTER TABLE personal_info ADD COLUMN introduction TEXT NOT NULL DEFAULT ''`)
	return err
}

func (s *Store) Create(ctx context.Context, person model.Person) error {
	result, err := s.db.ExecContext(ctx, `
INSERT OR IGNORE INTO personal_info (userid, name, email, phone, introduction)
VALUES (?, ?, ?, ?, ?)`,
		person.UserID, person.Name, person.Email, person.Phone, person.Introduction,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDuplicate
	}

	return nil
}

func (s *Store) Get(ctx context.Context, userid string) (model.Person, error) {
	var person model.Person
	err := s.db.QueryRowContext(ctx, `
SELECT userid, name, email, phone, introduction
FROM personal_info
WHERE userid = ?`, userid).Scan(&person.UserID, &person.Name, &person.Email, &person.Phone, &person.Introduction)
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
SET name = ?, email = ?, phone = ?, introduction = ?
WHERE userid = ?`,
		person.Name, person.Email, person.Phone, person.Introduction, person.UserID,
	)
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
