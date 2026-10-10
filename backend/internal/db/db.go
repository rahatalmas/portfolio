// Package db opens the MySQL connection pool and applies the embedded SQL
// migrations in internal/db/migrations.
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Open connects to MySQL and verifies the connection.
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := conn.PingContext(pingCtx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return conn, nil
}

// Migrate applies every migrations/NNNN_name.sql file whose version is not yet
// recorded in schema_migrations, in ascending order. Files hold plain SQL
// statements separated by a semicolon at the end of a line. It returns the
// number of migrations applied.
func Migrate(ctx context.Context, conn *sql.DB) (int, error) {
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INT       NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (version)
	) ENGINE=InnoDB`); err != nil {
		return 0, fmt.Errorf("create schema_migrations: %w", err)
	}

	applied := map[int]bool{}
	rows, err := conn.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return 0, err
		}
		applied[v] = true
	}
	rows.Close()

	files, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return 0, err
	}
	sort.Strings(files)

	n := 0
	for _, name := range files {
		var version int
		base := strings.TrimPrefix(name, "migrations/")
		if _, err := fmt.Sscanf(base, "%d_", &version); err != nil {
			return n, fmt.Errorf("migration %q: name must start with a number", base)
		}
		if applied[version] {
			continue
		}
		body, err := migrationFS.ReadFile(name)
		if err != nil {
			return n, err
		}
		for _, stmt := range splitStatements(string(body)) {
			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				return n, fmt.Errorf("migration %s: %w\n%s", base, err, stmt)
			}
		}
		if _, err := conn.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// splitStatements splits a migration file on semicolons that end a line,
// dropping comment-only lines and blank statements.
func splitStatements(sqlText string) []string {
	var out []string
	var cur strings.Builder
	for _, line := range strings.Split(sqlText, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSpace(cur.String())
			out = append(out, strings.TrimSuffix(stmt, ";"))
			cur.Reset()
		}
	}
	if rest := strings.TrimSpace(cur.String()); rest != "" {
		out = append(out, rest)
	}
	return out
}
