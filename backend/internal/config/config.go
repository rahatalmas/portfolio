// Package config loads runtime configuration from the environment, with an
// optional .env file in the working directory for local development.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is everything the server needs to start.
type Config struct {
	Addr          string   // listen address, e.g. ":8080"
	DSN           string   // MySQL DSN for go-sql-driver/mysql
	AdminUsername string   // seeded on first run when no admin user exists
	AdminPassword string   // seeded on first run; generated and logged if empty
	SessionSecret string   // HMAC key for admin session cookies
	UploadDir     string   // where media uploads are stored
	FrontendDir   string   // built React app to serve at /; skipped if missing
	CORSOrigins   []string // allowed origins for the public API; empty = same-origin only
	SecureCookies bool     // set the Secure flag on cookies (enable behind HTTPS)
	SeedDemo      bool     // insert demo content when the content tables are empty
}

// Load reads .env (if present) then the process environment.
func Load() (Config, error) {
	_ = loadDotEnv(".env")

	c := Config{
		Addr:          env("ADDR", ":8080"),
		DSN:           os.Getenv("DB_DSN"),
		AdminUsername: env("ADMIN_USERNAME", "admin"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		UploadDir:     env("UPLOAD_DIR", "./uploads"),
		FrontendDir:   env("FRONTEND_DIR", "../frontend/dist"),
		SecureCookies: envBool("SECURE_COOKIES", false),
		SeedDemo:      envBool("SEED_DEMO", true),
	}
	if o := strings.TrimSpace(os.Getenv("CORS_ORIGINS")); o != "" {
		for _, s := range strings.Split(o, ",") {
			if s = strings.TrimSpace(s); s != "" {
				c.CORSOrigins = append(c.CORSOrigins, s)
			}
		}
	}

	if c.DSN == "" {
		// Build a DSN from the individual DB_* variables.
		host := env("DB_HOST", "127.0.0.1")
		port := env("DB_PORT", "3306")
		user := env("DB_USER", "portfolio")
		pass := os.Getenv("DB_PASSWORD")
		name := env("DB_NAME", "portfolio")
		c.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	}
	// Options the server relies on: timestamps as time.Time, utf8mb4 everywhere.
	if !strings.Contains(c.DSN, "?") {
		c.DSN += "?"
	} else if !strings.HasSuffix(c.DSN, "&") && !strings.HasSuffix(c.DSN, "?") {
		c.DSN += "&"
	}
	c.DSN += "parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	if len(c.SessionSecret) > 0 && len(c.SessionSecret) < 16 {
		return c, errors.New("SESSION_SECRET must be at least 16 characters")
	}
	return c, nil
}

func env(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

// loadDotEnv sets variables from a KEY=VALUE file without overriding ones
// already present in the environment. Lines starting with # are comments;
// values may be wrapped in single or double quotes.
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if len(v) >= 2 && (v[0] == '"' && v[len(v)-1] == '"' || v[0] == '\'' && v[len(v)-1] == '\'') {
			v = v[1 : len(v)-1]
		}
		if _, exists := os.LookupEnv(k); !exists {
			os.Setenv(k, v)
		}
	}
	return sc.Err()
}
