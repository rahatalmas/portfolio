// Package api serves the public JSON API under /api/v1. All endpoints are
// read-only except POST /api/v1/messages (the contact form). See docs/API.md.
package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/schema"
	"portfolio/backend/internal/store"
)

// Handler holds the API's dependencies.
type Handler struct {
	store   *store.Store
	log     *slog.Logger
	limiter *rateLimiter
}

// New constructs the API handler.
func New(st *store.Store, log *slog.Logger) *Handler {
	return &Handler{store: st, log: log, limiter: newRateLimiter(5, time.Minute)}
}

// Routes registers every endpoint on mux.
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", h.health)
	mux.HandleFunc("GET /api/v1/portfolio", h.portfolio)
	mux.HandleFunc("GET /api/v1/site", h.site)
	mux.HandleFunc("GET /api/v1/education", h.education)
	mux.HandleFunc("GET /api/v1/skills", h.skills)
	mux.HandleFunc("GET /api/v1/profiles", h.profiles)
	mux.HandleFunc("GET /api/v1/experience", h.experience)
	mux.HandleFunc("GET /api/v1/projects", h.projects)
	mux.HandleFunc("GET /api/v1/papers", h.papers)
	mux.HandleFunc("POST /api/v1/messages", h.createMessage)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, "no such endpoint")
	})
}

// cacheable marks a public GET response as briefly cacheable so a burst of
// visitors is served from the browser cache instead of a new request.
func cacheable(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "public, max-age=30")
}

func (h *Handler) fail(w http.ResponseWriter, what string, err error) {
	h.log.Error("api", "what", what, "err", err)
	httpx.Error(w, http.StatusInternalServerError, "could not load "+what)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) portfolio(w http.ResponseWriter, r *http.Request) {
	p, err := h.store.Portfolio(r.Context())
	if err != nil {
		h.fail(w, "portfolio", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, p)
}

func (h *Handler) site(w http.ResponseWriter, r *http.Request) {
	s, err := h.store.Settings(r.Context())
	if err != nil {
		h.fail(w, "site", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, schema.TypedSettings(s))
}

func (h *Handler) education(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Education(r.Context())
	if err != nil {
		h.fail(w, "education", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

func (h *Handler) skills(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Skills(r.Context())
	if err != nil {
		h.fail(w, "skills", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

func (h *Handler) profiles(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Profiles(r.Context())
	if err != nil {
		h.fail(w, "profiles", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

func (h *Handler) experience(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Experience(r.Context())
	if err != nil {
		h.fail(w, "experience", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

func (h *Handler) projects(w http.ResponseWriter, r *http.Request) {
	featured := r.URL.Query().Get("featured")
	v, err := h.store.Projects(r.Context(), featured == "1" || featured == "true")
	if err != nil {
		h.fail(w, "projects", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

func (h *Handler) papers(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Papers(r.Context())
	if err != nil {
		h.fail(w, "papers", err)
		return
	}
	cacheable(w)
	httpx.JSON(w, http.StatusOK, v)
}

// messageInput is the contact form payload. `website` is a honeypot: real
// users never see it, bots fill it, and we silently drop those.
type messageInput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Website string `json:"website"`
}

func (h *Handler) createMessage(w http.ResponseWriter, r *http.Request) {
	if !h.limiter.Allow(httpx.ClientIP(r)) {
		httpx.Error(w, http.StatusTooManyRequests, "too many messages, try again in a minute")
		return
	}

	var in messageInput
	body, err := io.ReadAll(io.LimitReader(r.Body, 64<<10))
	if err != nil || json.Unmarshal(body, &in) != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(in.Email)
	in.Message = strings.TrimSpace(in.Message)

	if problem := validateMessage(in); problem != "" {
		httpx.Error(w, http.StatusUnprocessableEntity, problem)
		return
	}
	if in.Website != "" {
		// Honeypot tripped: pretend it worked.
		httpx.JSON(w, http.StatusCreated, map[string]bool{"ok": true})
		return
	}

	id, err := h.store.CreateMessage(r.Context(), in.Name, in.Email, in.Message)
	if err != nil {
		h.fail(w, "message", err)
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]any{"ok": true, "id": id})
}

func validateMessage(in messageInput) string {
	switch {
	case in.Name == "":
		return "name is required"
	case utf8.RuneCountInString(in.Name) > 120:
		return "name is too long (max 120 characters)"
	case in.Email == "":
		return "email is required"
	case utf8.RuneCountInString(in.Email) > 254:
		return "email is too long"
	case in.Message == "":
		return "message is required"
	case utf8.RuneCountInString(in.Message) > 5000:
		return "message is too long (max 5000 characters)"
	}
	if _, err := mail.ParseAddress(in.Email); err != nil {
		return "email is not a valid address"
	}
	return ""
}
