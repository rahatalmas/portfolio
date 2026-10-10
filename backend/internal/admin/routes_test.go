package admin

import (
	"log/slog"
	"net/http"
	"os"
	"testing"

	"portfolio/backend/internal/media"
)

// ServeMux panics at registration time when two patterns are ambiguous, so
// simply registering the routes is a meaningful test.
func TestRoutesRegisterWithoutConflicts(t *testing.T) {
	lib, err := media.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	h, err := New(nil, NewSessions([]byte("0123456789abcdef0123456789abcdef"), false), lib, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	h.Routes(mux) // panics on conflict
}
