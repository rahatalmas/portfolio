package store

import (
	"context"

	"portfolio/backend/internal/models"
)

// --- admin users (admin.json) ----------------------------------------------

// user is the on-disk shape of an admin account.
type user struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// UserByUsername returns the admin account with that username.
func (s *Store) UserByUsername(ctx context.Context, username string) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Username == username {
			return models.User{ID: u.ID, Username: u.Username, PasswordHash: u.PasswordHash}, nil
		}
	}
	return models.User{}, ErrNotFound
}

// CreateUser adds an admin account.
func (s *Store) CreateUser(ctx context.Context, username, passwordHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var maxID int64
	for _, u := range s.users {
		if u.Username == username {
			return ErrDuplicate
		}
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	s.users = append(s.users, user{ID: maxID + 1, Username: username, PasswordHash: passwordHash})
	if err := s.saveJSON(adminFile, s.users); err != nil {
		s.users = s.users[:len(s.users)-1]
		return err
	}
	return nil
}

// UpdatePassword replaces an account's password hash.
func (s *Store) UpdatePassword(ctx context.Context, username, passwordHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.users {
		if s.users[i].Username == username {
			s.users[i].PasswordHash = passwordHash
			return s.saveJSON(adminFile, s.users)
		}
	}
	return ErrNotFound
}

// UserCount returns how many admin accounts exist.
func (s *Store) UserCount(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users), nil
}
