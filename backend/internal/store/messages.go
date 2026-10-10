package store

import (
	"context"
	"sort"
	"time"

	"portfolio/backend/internal/models"
)

// message is the on-disk shape of a contact-form submission (messages.json).
type message struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Body      string    `json:"body"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

func (m message) model() models.Message {
	return models.Message{
		ID: m.ID, Name: m.Name, Email: m.Email, Body: m.Body, Read: m.Read,
		CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// CreateMessage stores a contact-form submission and returns its id.
func (s *Store) CreateMessage(ctx context.Context, name, email, body string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var maxID int64
	for _, m := range s.messages {
		if m.ID > maxID {
			maxID = m.ID
		}
	}
	m := message{ID: maxID + 1, Name: name, Email: email, Body: body, CreatedAt: time.Now().UTC()}
	s.messages = append(s.messages, m)
	if err := s.saveJSON(messagesFile, s.messages); err != nil {
		s.messages = s.messages[:len(s.messages)-1]
		return 0, err
	}
	return m.ID, nil
}

// Messages returns submissions newest first, at most limit (default 200).
func (s *Store) Messages(ctx context.Context, limit int) ([]models.Message, error) {
	if limit <= 0 {
		limit = 200
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sorted := make([]message, len(s.messages))
	copy(sorted, s.messages)
	sort.SliceStable(sorted, func(i, j int) bool {
		if !sorted[i].CreatedAt.Equal(sorted[j].CreatedAt) {
			return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
		}
		return sorted[i].ID > sorted[j].ID
	})
	out := []models.Message{}
	for i, m := range sorted {
		if i >= limit {
			break
		}
		out = append(out, m.model())
	}
	return out, nil
}

// Message returns one submission.
func (s *Store) Message(ctx context.Context, id int64) (models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i := s.messageIndex(id); i >= 0 {
		return s.messages[i].model(), nil
	}
	return models.Message{}, ErrNotFound
}

// MarkMessageRead flags a submission as read or unread.
func (s *Store) MarkMessageRead(ctx context.Context, id int64, read bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.messageIndex(id)
	if i < 0 {
		return ErrNotFound
	}
	if s.messages[i].Read == read {
		return nil
	}
	s.messages[i].Read = read
	return s.saveJSON(messagesFile, s.messages)
}

// DeleteMessage removes a submission.
func (s *Store) DeleteMessage(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.messageIndex(id)
	if i < 0 {
		return ErrNotFound
	}
	s.messages = append(s.messages[:i:i], s.messages[i+1:]...)
	return s.saveJSON(messagesFile, s.messages)
}

// UnreadMessages counts submissions not yet read.
func (s *Store) UnreadMessages(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, m := range s.messages {
		if !m.Read {
			n++
		}
	}
	return n, nil
}

// messageIndex returns the position of id, or -1. Caller holds mu.
func (s *Store) messageIndex(id int64) int {
	for i, m := range s.messages {
		if m.ID == id {
			return i
		}
	}
	return -1
}

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
