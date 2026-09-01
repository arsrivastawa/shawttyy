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
	Get(shortURL string) (*models.URL, error)
	Delete(shortURL string) error
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

	s.urls[url.ShortURL] = url
	return nil
}

func (s *InMemoryStore) Get(shortURL string) (*models.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.urls[shortURL]
	if !ok {
		return nil, ErrNotFound
	}
	return url, nil
}

func (s *InMemoryStore) Delete(shortURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.urls[shortURL]; !ok {
		return ErrNotFound
	}
	delete(s.urls, shortURL)
	return nil
}
