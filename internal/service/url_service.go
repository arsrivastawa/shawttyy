package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/arsrivastawa/shawttyy/internal/cache"
	"github.com/arsrivastawa/shawttyy/internal/core/encoder"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/exceptions"
	"github.com/arsrivastawa/shawttyy/internal/storage"
	"github.com/arsrivastawa/shawttyy/models"
)

const (
	defaultExpiry = 5 * 365 * 24 * time.Hour
	maxAliasLen   = 11
)

type URLService struct {
	store storage.Store
	cache cache.Cache
	seq   *sequencer.Sequencer
}

func New(store storage.Store, cache cache.Cache, seq *sequencer.Sequencer) *URLService {
	return &URLService{store: store, cache: cache, seq: seq}
}

// Shorten creates a URL entry and returns it. When customAlias is empty a new
// generated short code is used; otherwise the alias is validated and checked
// for collisions.
func (s *URLService) Shorten(req *models.CreateURLRequest) (*models.URL, error) {
	if req.OriginalURL == "" {
		return nil, exceptions.ErrInvalidOriginalURL
	}

	longURL := req.OriginalURL
	shortCode := req.CustomAlias
	userID := req.UserID
	isCustom := shortCode != ""

	{
		url, err := s.store.GetByLongURLAndUserID(longURL, userID)

		if err == nil {
			return url, nil
		}
	}

	// fmt.Print(" user idfrom service", userID)

	if isCustom {
		if !validAlias(shortCode) {
			return nil, exceptions.ErrInvalidCustomAlias
		}
		if _, err := s.store.Get(shortCode); err == nil {
			return nil, exceptions.ErrCustomAliasTaken
		} else if !errors.Is(err, exceptions.ErrNotFound) {
			return nil, err
		}
	} else {
		shortCode = encoder.Base62Encode(s.seq.Next())
	}

	expiry := req.ExpirationTime
	if expiry.IsZero() {
		d := time.Now().Add(defaultExpiry)
		expiry = d
	}

	url := &models.URL{
		ID:          s.seq.Next(),
		UserID:      userID,
		ShortCode:   shortCode,
		OriginalURL: longURL,
		IsCustom:    isCustom,
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   &expiry,
	}

	s.cache.Set(shortCode, req.OriginalURL)

	if err := s.store.Save(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (s *URLService) Resolve(shortCode string) (string, error) {
	start := time.Now()

	val, err := s.cache.Get(shortCode)
	if err == nil {
		duration := time.Since(start)
		fmt.Printf("[CACHE HIT] Resolved /%s in %v\n", shortCode, duration)

		return val, nil
	}

	dbStart := time.Now()

	url, err := s.store.Get(shortCode)
	if err != nil {
		if errors.Is(err, exceptions.ErrNotFound) {
			return "", exceptions.ErrShortURLNotFound
		}
		return "", err
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return "", exceptions.ErrURLExpired
	}

	s.cache.Set(url.ShortCode, url.OriginalURL)

	totalDuration := time.Since(start)
	dbDuration := time.Since(dbStart)

	fmt.Printf("[DB HIT] Resolved /%s in %v (DB Query took %v)\n", shortCode, totalDuration, dbDuration)

	return url.OriginalURL, nil
}

func (s *URLService) Delete(shortCode string, userID string) error {

	{
		url, err := s.store.Get(shortCode)
		if err != nil {
			if errors.Is(err, exceptions.ErrNotFound) {
				return exceptions.ErrShortURLNotFound
			}
			return err
		}

		if url.UserID != userID {
			return exceptions.ErrUnauthorized
		}
	}

	err := s.store.Delete(shortCode, userID)
	if errors.Is(err, exceptions.ErrNotFound) {
		return exceptions.ErrShortURLNotFound
	}
	s.cache.Del(shortCode)
	return err
}

func validAlias(alias string) bool {
	if len(alias) > maxAliasLen {
		return false
	}
	for _, r := range alias {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return len(alias) > 0
}
