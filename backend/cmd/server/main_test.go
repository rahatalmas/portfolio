package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"portfolio/backend/internal/admin"
	"portfolio/backend/internal/api"
	"portfolio/backend/internal/media"
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
	mux := newMux(api.New(nil, log), adminH, lib)

	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}
	if rec := get("/"); rec.Code != http.StatusFound || rec.Header().Get("Location") != "/admin" {
		t.Fatalf("root should redirect to /admin: %d %q", rec.Code, rec.Header().Get("Location"))
	}
	if rec := get("/api/v1/nope"); rec.Code != 404 {
		t.Fatalf("api fallback: %d", rec.Code)
	}
	if rec := get("/admin/projects"); rec.Code != http.StatusSeeOther {
		t.Fatalf("admin without session should redirect to login: %d", rec.Code)
	}
	if rec := get("/uploads/"); rec.Code != 404 {
		t.Fatalf("upload directory listing should be 404: %d", rec.Code)
	}
	if rec := get("/some/site/path"); rec.Code != 404 {
		t.Fatalf("unknown paths are 404 (the site is not served here): %d", rec.Code)
	}
}
