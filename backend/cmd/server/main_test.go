package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"portfolio/backend/internal/admin"
	"portfolio/backend/internal/api"
	"portfolio/backend/internal/media"
	"portfolio/backend/internal/site"
)

// Registering the complete route table must not panic (ServeMux rejects
// ambiguous patterns at registration), and the top-level routes behave.
func TestNewMuxRegisters(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	lib, err := media.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	adminH, err := admin.New(nil, admin.NewSessions([]byte("0123456789abcdef0123456789abcdef"), false), lib, log)
	if err != nil {
		t.Fatal(err)
	}
	apiH := api.New(nil, log)

	get := func(mux http.Handler, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}
	common := func(mux http.Handler) {
		t.Helper()
		if rec := get(mux, "/api/v1/nope"); rec.Code != 404 {
			t.Fatalf("api fallback: %d", rec.Code)
		}
		if rec := get(mux, "/admin/projects"); rec.Code != http.StatusSeeOther {
			t.Fatalf("admin without session should redirect to login: %d", rec.Code)
		}
		if rec := get(mux, "/uploads/"); rec.Code != 404 {
			t.Fatalf("upload directory listing should be 404: %d", rec.Code)
		}
	}

	// Without a built site: root points at the admin, everything else is 404.
	mux := newMux(apiH, adminH, lib, nil)
	common(mux)
	if rec := get(mux, "/"); rec.Code != http.StatusFound || rec.Header().Get("Location") != "/admin" {
		t.Fatalf("root should redirect to /admin: %d %q", rec.Code, rec.Header().Get("Location"))
	}
	if rec := get(mux, "/some/site/path"); rec.Code != 404 {
		t.Fatalf("unknown paths are 404 without a site: %d", rec.Code)
	}

	// With a built site: root and client routes serve index.html; the API,
	// admin and uploads still win over the catch-all.
	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("<!doctype html>site"), 0o644); err != nil {
		t.Fatal(err)
	}
	siteH, err := site.Open(dist)
	if err != nil || siteH == nil {
		t.Fatalf("site.Open: %v %v", siteH, err)
	}
	mux = newMux(apiH, adminH, lib, siteH)
	common(mux)
	if rec := get(mux, "/"); rec.Code != 200 || rec.Body.String() != "<!doctype html>site" {
		t.Fatalf("root should serve index.html: %d %q", rec.Code, rec.Body.String())
	}
	if rec := get(mux, "/some/site/path"); rec.Code != 200 {
		t.Fatalf("client routes fall back to index.html: %d", rec.Code)
	}
}
