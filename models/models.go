package models

import (
	"time"
)

type URL struct {
	ID          int64     `json:"id"`
	IsCustom    bool      `json:"is_custom"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateURLRequest struct {
	APIKey         string    `json:"api_key"`
	OriginalURL    string    `json:"original_url"`
	CustomShortURL string    `json:"custom_short_url,omitempty"`
	ExpirationTime time.Time `json:"expiration_time"`
}

type CreateURLResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}
