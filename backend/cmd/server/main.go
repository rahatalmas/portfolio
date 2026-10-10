// Command server runs the whole portfolio from one process: the built React
// site at /, the public JSON API, the HTML admin, and uploaded media.
//
// Content and the admin account are JSON files in DATA_DIR; there
// is no database. Configuration comes from the environment or a .env file;
// see backend/.env.example and backend/README.md.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"

	"portfolio/backend/internal/admin"
	"portfolio/backend/internal/api"
	"portfolio/backend/internal/config"
	"portfolio/backend/internal/httpx"
	"portfolio/backend/internal/media"
	"portfolio/backend/internal/site"
	"portfolio/backend/internal/store"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(log); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// --- data --------------------------------------------------------------
	st, err := store.Open(cfg.DataDir)
	if err != nil {
		return err
	}
	log.Info("data loaded", "dir", st.Dir())
	if err := ensureAdminUser(ctx, st, cfg, log); err != nil {
		return err
	}

	// --- services ----------------------------------------------------------
	secret := []byte(cfg.SessionSecret)
	if len(secret) == 0 {
		secret = randomBytes(32)
		log.Warn("SESSION_SECRET not set: using a random one, admin sessions reset on restart")
	}
	sessions := admin.NewSessions(secret, cfg.SecureCookies)

	lib, err := media.New(cfg.UploadDir)
	if err != nil {
		return err
	}

	adminH, err := admin.New(st, sessions, lib, log)
	if err != nil {
		return err
	}
	apiH := api.New(st, log)

	siteH, err := site.Open(cfg.FrontendDir)
	if err != nil {
		return err
	}
	if siteH == nil {
		log.Warn("no built site found; / redirects to /admin (set FRONTEND_DIR or run `npm run build` in frontend/)", "dir", cfg.FrontendDir)
	}

	// --- routes ------------------------------------------------------------
	mux := newMux(apiH, adminH, lib, siteH)

	handler := httpx.Chain(mux,
		httpx.Recover(log),
		httpx.Logger(log),
		httpx.SecurityHeaders(),
		httpx.CORS(cfg.CORSOrigins),
		httpx.MaxBody(media.MaxUploadBytes+1<<20),
	)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("listening", "addr", cfg.Addr, "admin", "http://localhost"+portOf(cfg.Addr)+"/admin")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

// ensureAdminUser creates the first admin account when none exists, using
// ADMIN_USERNAME/ADMIN_PASSWORD or a generated password printed once.
func ensureAdminUser(ctx context.Context, st *store.Store, cfg config.Config, log *slog.Logger) error {
	n, err := st.UserCount(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	password := cfg.AdminPassword
	generated := false
	if password == "" {
		password = hex.EncodeToString(randomBytes(9))
		generated = true
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := st.CreateUser(ctx, cfg.AdminUsername, string(hash)); err != nil {
		return err
	}
	if generated {
		log.Warn("created admin user with a generated password — change it after logging in",
			"username", cfg.AdminUsername, "password", password)
	} else {
		log.Info("created admin user", "username", cfg.AdminUsername)
	}
	return nil
}

// newMux assembles the route table: API, admin, uploads, and the React site
// on everything else. When no built site is available (siteH == nil) the
// root redirects to the admin and other paths are 404.
func newMux(apiH *api.Handler, adminH *admin.Handler, lib *media.Library, siteH *site.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	apiH.Routes(mux)
	adminH.Routes(mux)
	mux.Handle("GET "+media.URLPrefix, lib.Handler())
	if siteH != nil {
		mux.Handle("/", siteH)
		return mux
	}
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusFound)
	})
	return mux
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func portOf(addr string) string {
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		return addr[i:]
	}
	return ""
}
