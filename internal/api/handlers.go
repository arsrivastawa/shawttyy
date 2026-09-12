package api

import (
	"encoding/json"
	"errors"
	"net/http"

	// "net/http/httputil"

	"github.com/arsrivastawa/shawttyy/internal/core/ip"
	"github.com/arsrivastawa/shawttyy/internal/exceptions"

	"github.com/arsrivastawa/shawttyy/internal/service"
	"github.com/arsrivastawa/shawttyy/models"
)

type Handler struct {
	svc     *service.URLService
	baseURL string
}

func NewHandler(svc *service.URLService, baseURL string) *Handler {
	return &Handler{svc: svc, baseURL: baseURL}
}

// Shorten handles POST /shorten.
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req models.CreateURLRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	IP, err := ip.ExtractIPFromRequest(r.RemoteAddr)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "could not extract IP address")
		return
	}

	hash := ip.IPEncoder(IP)
	req.UserID = hash

	url, err := h.svc.Shorten(&req)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrInvalidOriginalURL):
			h.writeError(w, http.StatusBadRequest, exceptions.ErrInvalidOriginalURL.Error())
		case errors.Is(err, exceptions.ErrInvalidCustomAlias):
			h.writeError(w, http.StatusBadRequest, exceptions.ErrInvalidCustomAlias.Error())
		case errors.Is(err, exceptions.ErrCustomAliasTaken):
			h.writeError(w, http.StatusConflict, exceptions.ErrCustomAliasTaken.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, "could not shorten url")
			// println("Error shortening URL:", err.Error())
		}
		return
	}

	h.writeJSON(w, http.StatusCreated, models.CreateURLResponse{
		ShortURL:  h.baseURL + "/" + url.ShortCode,
		ShortCode: url.ShortCode,
	})
}

// Redirect handles GET /{short_code}.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	url, err := h.svc.Resolve(shortCode)
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrNotFound):
			h.writeError(w, http.StatusNotFound, exceptions.ErrShortURLNotFound.Error())
		case errors.Is(err, exceptions.ErrURLExpired):
			h.writeError(w, http.StatusNotFound, exceptions.ErrURLExpired.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, exceptions.ErrCanNotResolveURL.Error())
		}
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

// Delete handles DELETE /{short_code}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	if err := h.svc.Delete(shortCode); err != nil {
		if errors.Is(err, exceptions.ErrNotFound) || errors.Is(err, exceptions.ErrShortURLNotFound) {
			h.writeError(w, http.StatusNotFound, exceptions.ErrShortURLNotFound.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, exceptions.ErrCanNotResolveURL.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"message": "URL Removed"})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}
