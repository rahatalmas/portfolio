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
	// Unknown paths get the bundle (which renders its not-found page) under a
	// 404 status, not a 200 that would claim the URL exists.
	if rec := get(h, "/some/random/uri"); rec.Code != 404 || rec.Body.String() == "" ||
		rec.Header().Get("Cache-Control") != "no-cache" || rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("unknown path: %d %q %q", rec.Code, rec.Header().Get("Cache-Control"), rec.Header().Get("Content-Type"))
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
	if rec := get(h, "/assets/"); rec.Code != 404 || rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatalf("directory should fall back to the not-found page, not list: %d", rec.Code)
	}
	if rec := get(h, "/../index.html"); rec.Code != 200 {
		t.Fatalf("traversal is cleaned: %d", rec.Code)
	}
	// HEAD on an unknown path: the 404 and the length, no body.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/nope", nil))
	if rec.Code != 404 || rec.Body.Len() != 0 || rec.Header().Get("Content-Length") == "" {
		t.Fatalf("HEAD on unknown path: %d body=%d len=%q", rec.Code, rec.Body.Len(), rec.Header().Get("Content-Length"))
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST: %d", rec.Code)
	}
}

// The Route hook is what makes a client-side route return 200 while every
// other invented path still returns 404 with the not-found page.
func TestRouteHook(t *testing.T) {
	h, err := Open(build(t))
	if err != nil || h == nil {
		t.Fatalf("open: %v %v", h, err)
	}
	var asked []string
	h.Route = func(p string) bool {
		asked = append(asked, p)
		return p == "/projects/fast-rag"
	}

	if rec := get(h, "/projects/fast-rag"); rec.Code != 200 {
		t.Fatalf("a known route should be 200: %d", rec.Code)
	}
	// Route is asked the cleaned path, so a trailing slash or a dot segment
	// resolves to the same answer.
	if rec := get(h, "/projects/fast-rag/"); rec.Code != 200 {
		t.Fatalf("trailing slash should be 200: %d", rec.Code)
	}
	if rec := get(h, "/projects/nope"); rec.Code != 404 {
		t.Fatalf("an unknown slug should be 404: %d", rec.Code)
	}
	if rec := get(h, "/"); rec.Code != 200 {
		t.Fatalf("root does not need the hook: %d", rec.Code)
	}
	for _, p := range asked {
		if p != "/projects/fast-rag" && p != "/projects/nope" {
			t.Fatalf("hook asked about an unexpected path: %q", p)
		}
	}
	// A missing file is still a plain 404: the hook is never consulted for
	// something that looks like a file.
	if rec := get(h, "/projects/fast-rag.png"); rec.Code != 404 {
		t.Fatalf("missing file: %d", rec.Code)
	}
}
