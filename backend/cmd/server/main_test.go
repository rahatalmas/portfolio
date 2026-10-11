package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
		if rec := get(mux, "/uploads/"); rec.Code != 404 {
			t.Fatalf("upload directory listing should be 404: %d", rec.Code)
		}
	}

	// Without a built site: the root explains itself, everything else is 404.
	mux := newPublicMux(apiH, lib, nil, "http://localhost:8081/admin")
	common(mux)
	if rec := get(mux, "/"); rec.Code != 404 {
		t.Fatalf("root without a built site should be 404: %d", rec.Code)
	}
	if rec := get(mux, "/some/site/path"); rec.Code != 404 {
		t.Fatalf("unknown paths are 404 without a site: %d", rec.Code)
	}

	// With a built site: the root serves index.html and the API and uploads
	// still win over the catch-all.
	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("<!doctype html>site"), 0o644); err != nil {
		t.Fatal(err)
	}
	siteH, err := site.Open(dist)
	if err != nil || siteH == nil {
		t.Fatalf("site.Open: %v %v", siteH, err)
	}
	mux = newPublicMux(apiH, lib, siteH, "http://localhost:8081/admin")
	common(mux)
	if rec := get(mux, "/"); rec.Code != 200 || rec.Body.String() != "<!doctype html>site" {
		t.Fatalf("root should serve index.html: %d %q", rec.Code, rec.Body.String())
	}
	// Unknown paths still serve the bundle, so it can render its not-found
	// page, but under a 404 status.
	if rec := get(mux, "/some/site/path"); rec.Code != 404 || rec.Body.String() != "<!doctype html>site" {
		t.Fatalf("unknown path should be the bundle with a 404: %d %q", rec.Code, rec.Body.String())
	}
}

// The admin must not be reachable on the public listener at all: that is what
// keeps it off the domain, no matter which headers a caller sends.
func TestAdminIsOnlyOnTheAdminMux(t *testing.T) {
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

	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("<!doctype html>site"), 0o644); err != nil {
		t.Fatal(err)
	}
	siteH, err := site.Open(dist)
	if err != nil || siteH == nil {
		t.Fatalf("site.Open: %v %v", siteH, err)
	}

	pub := newPublicMux(apiH, lib, siteH, "http://localhost:8081/admin")
	for _, path := range []string{"/admin", "/admin/login", "/admin/projects", "/admin/assets/admin.css"} {
		for _, host := range []string{"example.com", "localhost"} {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Host = host
			rec := httptest.NewRecorder()
			pub.ServeHTTP(rec, req)
			// The site handler answers instead: index.html for client routes,
			// 404 for a missing asset. What matters is that no admin handler
			// runs, so there is never a redirect to the login form nor any
			// admin markup in the body.
			body := rec.Body.String()
			if rec.Code == http.StatusSeeOther || strings.Contains(body, "/admin/login") || strings.Contains(body, "csrf") {
				t.Fatalf("public mux served the admin for %s (host %s): %d %q", path, host, rec.Code, body)
			}
		}
	}

	adm := newAdminMux(adminH, lib)
	rec := httptest.NewRecorder()
	adm.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/projects", nil))
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("admin mux should redirect an anonymous caller to the login form: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	adm.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/login", nil))
	if rec.Code != 200 {
		t.Fatalf("admin mux should serve the login form: %d", rec.Code)
	}
}

// isLoopback decides whether the startup warning fires, so pin its answers.
func TestIsLoopback(t *testing.T) {
	for addr, want := range map[string]bool{
		"127.0.0.1:8081": true,
		"localhost:8081": true,
		"[::1]:8081":     true,
		"127.0.0.1":      true,
		":8081":          false,
		"0.0.0.0:8081":   false,
		"192.168.1.9:80": false,
	} {
		if got := isLoopback(addr); got != want {
			t.Errorf("isLoopback(%q) = %v, want %v", addr, got, want)
		}
	}
}
