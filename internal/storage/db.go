package storage

import (
	"database/sql"
	"errors"

	"github.com/arsrivastawa/shawttyy/internal/exceptions"
	"github.com/arsrivastawa/shawttyy/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Save(url *models.URL) error {
	query := `
        INSERT INTO urls (id, user_id, short_code, original_url, is_custom, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := s.db.Exec(query, url.ID, url.UserID, url.ShortCode, url.OriginalURL, url.IsCustom, url.CreatedAt, url.ExpiresAt)
	return err
}

func (s *PostgresStore) Get(shortCode string) (*models.URL, error) {
	query := `
        SELECT id, user_id, short_code, original_url, is_custom, created_at, expires_at 
        FROM urls 
        WHERE short_code = $1
    `
	row := s.db.QueryRow(query, shortCode)

	var url models.URL
	err := row.Scan(
		&url.ID, &url.UserID, &url.ShortCode, &url.OriginalURL,
		&url.IsCustom, &url.CreatedAt, &url.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}

	return &url, nil
}

func (s *PostgresStore) Delete(shortCode string, userID string) error {
	query := `DELETE FROM urls WHERE short_code = $1 AND user_id = $2`
	result, err := s.db.Exec(query, shortCode, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return exceptions.ErrNotFound
	}

	return nil
}

func (s *PostgresStore) GetByLongURLAndUserID(longURL, userID string) (*models.URL, error) {
	query := `SELECT * FROM urls WHERE original_url = $1 AND user_id = $2`

	row := s.db.QueryRow(query, longURL, userID)

	var url models.URL
	err := row.Scan(
		&url.ID, &url.UserID, &url.ShortCode, &url.OriginalURL,
		&url.IsCustom, &url.CreatedAt, &url.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exceptions.ErrNotFound
		}
		return nil, err
	}

	return &url, nil
}
