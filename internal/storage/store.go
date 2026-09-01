package storage

import (
	"errors"
	"sync"

	"github.com/arsrivastawa/shawttyy/models"
)

var (
	ErrNotFound = errors.New("URL not found")
)

// Store is the persistence contract for URL mappings. Swap the in-memory
// implementation for Postgres in a later milestone without touching callers.
type Store interface {
	Save(url *models.URL) error
	Get(shortCode string) (*models.URL, error)
	Delete(shortCode string) error
}

type InMemoryStore struct {
	mu   sync.RWMutex
	urls map[string]*models.URL
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{urls: make(map[string]*models.URL)}
}

func (s *InMemoryStore) Save(url *models.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls[url.ShortCode] = url
	return nil
}

func (s *InMemoryStore) Get(shortCode string) (*models.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[shortCode]
	if !ok {
		return nil, ErrNotFound
	}
	return url, nil
}

func (s *InMemoryStore) Delete(shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.urls[shortCode]; !ok {
		return ErrNotFound
	}
	delete(s.urls, shortCode)
	return nil
}
