package models

import (
	"time"
)

type URL struct {
	ID          int64      `json:"id"`
	UserID      string     `json:"user_id"`
	IsCustom    bool       `json:"is_custom"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateURLRequest struct {
	UserID         string    `json:"user_id"`
	OriginalURL    string    `json:"original_url"`
	CustomAlias    string    `json:"custom_alias,omitempty"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type CreateURLResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}
