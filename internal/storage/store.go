package storage

import (
	"github.com/arsrivastawa/shawttyy/models"
)

// Store is the persistence contract for URL mappings. Swap the in-memory
// implementation for Postgres in a later milestone without touching callers.
type Store interface {
	Save(url *models.URL) error
	Get(shortCode string) (*models.URL, error)
	Delete(shortCode string, userID string) error
	GetByLongURLAndUserID(longURL, userID string) (*models.URL, error)
}
