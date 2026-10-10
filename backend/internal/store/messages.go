package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"portfolio/backend/internal/models"
)

// CreateMessage stores a contact-form submission.
func (s *Store) CreateMessage(ctx context.Context, name, email, body string) (int64, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO messages (name, email, body) VALUES (?, ?, ?)", name, email, body)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Messages returns submissions newest first.
func (s *Store) Messages(ctx context.Context, limit int) ([]models.Message, error) {
	if limit <= 0 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, email, body, is_read, created_at FROM messages ORDER BY created_at DESC, id DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Message{}
	for rows.Next() {
		m, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// Message returns one submission.
func (s *Store) Message(ctx context.Context, id int64) (models.Message, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, email, body, is_read, created_at FROM messages WHERE id = ?", id)
	if err != nil {
		return models.Message{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return models.Message{}, ErrNotFound
	}
	return scanMessage(rows)
}

// MarkMessageRead flags a submission as read or unread.
func (s *Store) MarkMessageRead(ctx context.Context, id int64, read bool) error {
	_, err := s.db.ExecContext(ctx, "UPDATE messages SET is_read = ? WHERE id = ?", read, id)
	return err
}

// DeleteMessage removes a submission.
func (s *Store) DeleteMessage(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// UnreadMessages counts submissions not yet read.
func (s *Store) UnreadMessages(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM messages WHERE is_read = 0").Scan(&n)
	return n, err
}

func scanMessage(rows *sql.Rows) (models.Message, error) {
	var m models.Message
	var created time.Time
	if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Body, &m.Read, &created); err != nil {
		return m, err
	}
	m.CreatedAt = created.UTC().Format(time.RFC3339)
	return m, nil
}

// --- admin users ---------------------------------------------------------

// UserByUsername returns the admin account with that username.
func (s *Store) UserByUsername(ctx context.Context, username string) (models.User, error) {
	var u models.User
	err := s.db.QueryRowContext(ctx, "SELECT id, username, password_hash FROM admin_users WHERE username = ?", username).
		Scan(&u.ID, &u.Username, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

// CreateUser adds an admin account.
func (s *Store) CreateUser(ctx context.Context, username, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO admin_users (username, password_hash) VALUES (?, ?)", username, passwordHash)
	return wrapDuplicate(err)
}

// UpdatePassword replaces an account's password hash.
func (s *Store) UpdatePassword(ctx context.Context, username, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE admin_users SET password_hash = ? WHERE username = ?", passwordHash, username)
	return err
}

// UserCount returns how many admin accounts exist.
func (s *Store) UserCount(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admin_users").Scan(&n)
	return n, err
}
