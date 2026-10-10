// Package api serves the public, read-only JSON API under /api/v1. See
// docs/API.md.
package api

import (
	"log/slog"
	"net/http"

	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/schema"
	"portfolio/backend/internal/store"
)

// Handler holds the API's dependencies.
type Handler struct {
	store *store.Store
	log   *slog.Logger
}

// New constructs the API handler.
func New(st *store.Store, log *slog.Logger) *Handler {
	return &Handler{store: st, log: log}
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
