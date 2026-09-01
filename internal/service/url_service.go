package service

import (
	"errors"
	"time"

	"github.com/arsrivastawa/shawttyy/internal/core/encoder"
	"github.com/arsrivastawa/shawttyy/internal/core/sequencer"
	"github.com/arsrivastawa/shawttyy/internal/storage"
	"github.com/arsrivastawa/shawttyy/models"
)

const (
	defaultExpiry = 5 * 365 * 24 * time.Hour
	maxAliasLen   = 11
)

var (
	ErrInvalidURL   = errors.New("service: original_url is required")
	ErrInvalidAlias = errors.New("service: custom_alias must be alphanumeric and <= 11 chars")
	ErrAliasTaken   = errors.New("service: custom_alias is already in use")
	ErrNotFound     = errors.New("service: short_url not found")
	ErrExpired      = errors.New("service: short_url has expired")
)

type URLService struct {
	store storage.Store
	seq   *sequencer.Sequencer
}

func New(store storage.Store, seq *sequencer.Sequencer) *URLService {
	return &URLService{store: store, seq: seq}
}

// Shorten creates a URL entry and returns it. When customAlias is empty a new
// generated short code is used; otherwise the alias is validated and checked
// for collisions.
func (s *URLService) Shorten(req *models.CreateURLRequest) (*models.URL, error) {
	if req.OriginalURL == "" {
		return nil, ErrInvalidURL
	}

	longURL := req.OriginalURL
	shortCode := req.CustomAlias
	isCustom := shortCode != ""

	if isCustom {
		if !validAlias(shortCode) {
			return nil, ErrInvalidAlias
		}
		if _, err := s.store.Get(shortCode); err == nil {
			return nil, ErrAliasTaken
		} else if !errors.Is(err, storage.ErrNotFound) {
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
		ShortCode:   shortCode,
		OriginalURL: longURL,
		IsCustom:    isCustom,
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   &expiry,
	}

	if err := s.store.Save(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (s *URLService) Resolve(shortCode string) (*models.URL, error) {
	url, err := s.store.Get(shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return nil, ErrExpired
	}

	return url, nil
}

func (s *URLService) Delete(shortCode string) error {
	err := s.store.Delete(shortCode)
	if errors.Is(err, storage.ErrNotFound) {
		return ErrNotFound
	}
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
