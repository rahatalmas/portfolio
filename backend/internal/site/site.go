// Package site serves the built React app (frontend/dist) from the Go server
// so the site, API, admin and uploads share one origin and one deployment.
package site

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Handler serves a Vite build directory. Requests for existing files are
// served as-is; every other GET falls back to index.html. The fallback carries
// a 404 status for anything but the root, because the site is one page and the
// bundle renders its own not-found page there: the status and the page agree,
// so crawlers are not told a missing URL is fine. Hashed files under /assets/
// are cached for a year; everything else (notably index.html) is revalidated.
type Handler struct {
	root  fs.FS
	files http.Handler
}

// Open returns a handler for dir, or (nil, nil) when dir does not contain a
// built site (no index.html) so the caller can run without the frontend.
func Open(dir string) (*Handler, error) {
	if dir == "" {
		return nil, nil
	}
	root := os.DirFS(dir)
	if _, err := fs.Stat(root, "index.html"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return &Handler{root: root, files: http.FileServerFS(root)}, nil
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	root := name == "" || name == "index.html"
	if name == "" {
		name = "index.html"
	}
	if name != "index.html" {
		if st, err := fs.Stat(h.root, name); err == nil && !st.IsDir() {
			if strings.HasPrefix(name, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else {
				w.Header().Set("Cache-Control", "no-cache")
			}
			h.files.ServeHTTP(w, r)
			return
		}
		// Looks like a file (has an extension) but is missing: a real 404 is
		// more useful than index.html with the wrong content type.
		if path.Ext(name) != "" {
			http.NotFound(w, r)
			return
		}
	}
	h.serveIndex(w, r, root)
}

// serveIndex writes index.html. ok serves it as the page it is; otherwise the
// same bytes go out with a 404 status for the bundle's not-found page.
func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request, ok bool) {
	f, err := h.root.Open("index.html")
	if err != nil {
		http.Error(w, "site unavailable", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		http.Error(w, "site unavailable", http.StatusInternalServerError)
		return
	}
	rs, isSeeker := f.(io.ReadSeeker)
	if !isSeeker {
		http.Error(w, "site unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if ok {
		http.ServeContent(w, r, "index.html", st.ModTime(), rs)
		return
	}
	// ServeContent cannot write a 404, and conditional requests and ranges are
	// meaningless for an error page, so send it plainly.
	body, err := io.ReadAll(rs)
	if err != nil {
		http.Error(w, "site unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusNotFound)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
}
