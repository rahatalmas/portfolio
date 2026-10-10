// Package media manages uploaded files (photo, CV, images) on disk. Files
// live flat in one directory and are served at /uploads/<name>.
package media

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// URLPrefix is where the server mounts the upload directory.
const URLPrefix = "/uploads/"

// MaxUploadBytes caps a single upload.
const MaxUploadBytes = 15 << 20 // 15 MiB

var allowedExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".svg": true,
	".pdf": true,
}

var imageExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true, ".svg": true}

var unsafeChars = regexp.MustCompile(`[^a-z0-9._-]+`)

// File describes one stored upload.
type File struct {
	Name    string
	Size    int64
	ModTime time.Time
	URL     string
	IsImage bool
}

// Library is an upload directory.
type Library struct {
	dir string
}

// New creates the directory if needed and returns a Library.
func New(dir string) (*Library, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &Library{dir: abs}, nil
}

// Dir returns the absolute upload directory.
func (l *Library) Dir() string { return l.dir }

// List returns stored files, newest first.
func (l *Library) List() ([]File, error) {
	entries, err := os.ReadDir(l.dir)
	if err != nil {
		return nil, err
	}
	out := make([]File, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, l.describe(e.Name(), info))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModTime.After(out[j].ModTime) })
	return out, nil
}

// Save stores r under a sanitised version of name, adding -2, -3… when the
// name is taken. It rejects disallowed extensions.
func (l *Library) Save(name string, r io.Reader) (File, error) {
	base := sanitize(name)
	ext := filepath.Ext(base)
	if !allowedExt[ext] {
		return File{}, fmt.Errorf("file type %q is not allowed", ext)
	}
	stem := strings.TrimSuffix(base, ext)

	final := base
	for i := 2; ; i++ {
		if _, err := os.Stat(filepath.Join(l.dir, final)); errors.Is(err, os.ErrNotExist) {
			break
		}
		final = fmt.Sprintf("%s-%d%s", stem, i, ext)
	}

	path := filepath.Join(l.dir, final)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return File{}, err
	}
	n, err := io.Copy(f, io.LimitReader(r, MaxUploadBytes+1))
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err == nil && n > MaxUploadBytes {
		err = fmt.Errorf("file is larger than %d MB", MaxUploadBytes>>20)
	}
	if err != nil {
		os.Remove(path)
		return File{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return l.describe(final, info), nil
}

// Delete removes a stored file by name.
func (l *Library) Delete(name string) error {
	if name == "" || name != filepath.Base(name) || strings.HasPrefix(name, ".") {
		return errors.New("invalid file name")
	}
	err := os.Remove(filepath.Join(l.dir, name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("file not found")
	}
	return err
}

// Handler serves the directory read-only at URLPrefix, without listings.
func (l *Library) Handler() http.Handler {
	fs := http.FileServer(http.Dir(l.dir))
	return http.StripPrefix(URLPrefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fs.ServeHTTP(w, r)
	}))
}

func (l *Library) describe(name string, info os.FileInfo) File {
	return File{
		Name:    name,
		Size:    info.Size(),
		ModTime: info.ModTime(),
		URL:     URLPrefix + name,
		IsImage: imageExt[strings.ToLower(filepath.Ext(name))],
	}
}

// sanitize turns an arbitrary upload name into a safe lowercase filename.
func sanitize(name string) string {
	name = strings.ToLower(filepath.Base(strings.ReplaceAll(name, "\\", "/")))
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	stem = strings.Trim(unsafeChars.ReplaceAllString(stem, "-"), "-.")
	if stem == "" {
		stem = "file"
	}
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	return stem + ext
}
