package site

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func build(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.MkdirAll(filepath.Join(dir, "assets"), 0o755))
	must(os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>site</title>"), 0o644))
	must(os.WriteFile(filepath.Join(dir, "assets", "app-abc123.js"), []byte("console.log(1)"), 0o644))
	must(os.WriteFile(filepath.Join(dir, "profile.jpg"), []byte("jpg"), 0o644))
	return dir
}

func get(h http.Handler, p string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, p, nil))
	return rec
}

func TestOpenMissingDir(t *testing.T) {
	h, err := Open(filepath.Join(t.TempDir(), "nope"))
	if err != nil || h != nil {
		t.Fatalf("missing dir should be (nil, nil): %v %v", h, err)
	}
	h, err = Open("")
	if err != nil || h != nil {
		t.Fatalf("empty dir should be (nil, nil): %v %v", h, err)
	}
}

func TestServe(t *testing.T) {
	h, err := Open(build(t))
	if err != nil || h == nil {
		t.Fatalf("open: %v %v", h, err)
	}

	if rec := get(h, "/"); rec.Code != 200 || rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("root: %d %q", rec.Code, rec.Header().Get("Content-Type"))
	}
	if rec := get(h, "/some/client/route"); rec.Code != 200 || rec.Body.String() == "" || rec.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("spa fallback: %d %q", rec.Code, rec.Header().Get("Cache-Control"))
	}
	if rec := get(h, "/assets/app-abc123.js"); rec.Code != 200 || rec.Header().Get("Cache-Control") != "public, max-age=31536000, immutable" {
		t.Fatalf("asset: %d %q", rec.Code, rec.Header().Get("Cache-Control"))
	}
	if rec := get(h, "/profile.jpg"); rec.Code != 200 || rec.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("public file: %d %q", rec.Code, rec.Header().Get("Cache-Control"))
	}
	if rec := get(h, "/missing.png"); rec.Code != 404 {
		t.Fatalf("missing file should be 404: %d", rec.Code)
	}
	if rec := get(h, "/assets/"); rec.Code != 200 || rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("directory should fall back to index, not list: %d", rec.Code)
	}
	if rec := get(h, "/../index.html"); rec.Code != 200 {
		t.Fatalf("traversal is cleaned: %d", rec.Code)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST: %d", rec.Code)
	}
}
