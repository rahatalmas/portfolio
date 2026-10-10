package admin

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"portfolio/backend/internal/media"
	"portfolio/backend/internal/schema"
)

//go:embed web
var webFS embed.FS

// view is the data every admin page receives. Page-specific values go in Data.
type view struct {
	Title     string
	User      string            // logged-in username; empty on the login page
	Handle    string            // navbar handle from settings, for the brand
	Active    string            // sidebar key to highlight
	Flash     string            // one-time notice
	FlashKind string            // "ok" | "error"
	Unread    int               // unread messages, for the sidebar badge
	Resources []schema.Resource // sidebar content links
	CSRF      string            // token for forms
	Data      any
}

// fieldView is one form input: its spec, current value and error.
type fieldView struct {
	F       schema.Field
	Value   string
	Error   string
	Checked map[string]bool // Tags kind: which options are on
}

// renderer holds the parsed templates, one set per page cloned from the
// shared layout so each page can define its own "content" block.
type renderer struct {
	pages map[string]*template.Template
	log   *slog.Logger
}

var funcs = template.FuncMap{
	"isImage": func(url string) bool {
		ext := strings.ToLower(path.Ext(url))
		return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif" || ext == ".svg"
	},
	"truncate": func(s string, n int) string {
		s = strings.ReplaceAll(s, "\n", ", ")
		r := []rune(s)
		if len(r) <= n {
			return s
		}
		return string(r[:n]) + "…"
	},
	"kb": func(n int64) string {
		if n >= 1<<20 {
			return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
		}
		return fmt.Sprintf("%d KB", (n+1023)/1024)
	},
	"inc": func(i int) int { return i + 1 },
}

func newRenderer(log *slog.Logger) (*renderer, error) {
	base, err := template.New("").Funcs(funcs).ParseFS(webFS, "web/templates/layout.html", "web/templates/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse layout: %w", err)
	}
	pageFiles, err := fs.Glob(webFS, "web/templates/pages/*.html")
	if err != nil {
		return nil, err
	}
	r := &renderer{pages: map[string]*template.Template{}, log: log}
	for _, file := range pageFiles {
		t, err := base.Clone()
		if err != nil {
			return nil, err
		}
		if _, err := t.ParseFS(webFS, file); err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		name := strings.TrimSuffix(path.Base(file), ".html")
		r.pages[name] = t
	}
	return r, nil
}

// render executes page with v. Errors are logged and turned into a 500.
func (r *renderer) render(w http.ResponseWriter, status int, page string, v view) {
	t, ok := r.pages[page]
	if !ok {
		r.log.Error("admin: unknown page", "page", page)
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := t.ExecuteTemplate(w, "layout", v); err != nil {
		r.log.Error("admin: render", "page", page, "err", err)
	}
}

// staticHandler serves the embedded admin assets (CSS).
func staticHandler() http.Handler {
	sub, err := fs.Sub(webFS, "web/static")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix("/admin/assets/", http.FileServer(http.FS(sub)))
}

// mediaOptions turns the library listing into datalist values.
func mediaOptions(files []media.File) []string {
	out := make([]string, len(files))
	for i, f := range files {
		out[i] = f.URL
	}
	return out
}
