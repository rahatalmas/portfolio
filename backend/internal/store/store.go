// Package store is the data-access layer. It has two halves:
//
//   - Typed reads for the public API (lists.go, settings.go), returning
//     models.* values.
//   - Generic, spec-driven CRUD for the admin (records.go), which works on
//     any schema.Resource so a new table needs no handler code.
//
// Plus messages (contact form) and admin users.
package store

import (
	"database/sql"
	"errors"
)

// ErrNotFound is returned when a row does not exist.
var ErrNotFound = errors.New("not found")

// Store wraps the connection pool.
type Store struct {
	db *sql.DB
}

// New returns a Store over an open pool.
func New(db *sql.DB) *Store { return &Store{db: db} }
