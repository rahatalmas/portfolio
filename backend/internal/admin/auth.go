package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	sessionCookie = "admin_session"
	flashCookie   = "admin_flash"
	sessionTTL    = 12 * time.Hour
)

// Sessions issues and verifies stateless, HMAC-signed session cookies. The
// cookie value is base64(username|expiry).base64(signature); nothing is
// stored server-side, so restarts don't log admins out (as long as the
// secret stays the same).
type Sessions struct {
	secret []byte
	secure bool
}

// NewSessions returns a session signer. secure sets the cookie Secure flag.
func NewSessions(secret []byte, secure bool) *Sessions {
	return &Sessions{secret: secret, secure: secure}
}

func (s *Sessions) sign(payload string) string {
	m := hmac.New(sha256.New, s.secret)
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// Issue builds a signed cookie value for username.
func (s *Sessions) Issue(username string) string {
	payload := username + "|" + strconv.FormatInt(time.Now().Add(sessionTTL).Unix(), 10)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + s.sign(payload)
}

// Verify checks a cookie value and returns the username if valid and unexpired.
func (s *Sessions) Verify(value string) (string, bool) {
	enc, sig, ok := strings.Cut(value, ".")
	if !ok {
		return "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(enc)
	if err != nil {
		return "", false
	}
	payload := string(raw)
	if !hmac.Equal([]byte(s.sign(payload)), []byte(sig)) {
		return "", false
	}
	username, expStr, ok := strings.Cut(payload, "|")
	if !ok || username == "" {
		return "", false
	}
	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return "", false
	}
	return username, true
}

// CSRF derives a per-session token from the cookie value. Forms embed it as
// _csrf and POST handlers compare it; an attacker's site cannot read the
// cookie, so it cannot forge the token.
func (s *Sessions) CSRF(cookieValue string) string {
	m := hmac.New(sha256.New, s.secret)
	m.Write([]byte("csrf:" + cookieValue))
	return hex.EncodeToString(m.Sum(nil))
}

// Set writes the session cookie.
func (s *Sessions) Set(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: s.Issue(username), Path: "/admin",
		HttpOnly: true, Secure: s.secure, SameSite: http.SameSiteLaxMode,
		MaxAge: int(sessionTTL.Seconds()),
	})
}

// Clear deletes the session cookie.
func (s *Sessions) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookie, Value: "", Path: "/admin", HttpOnly: true,
		Secure: s.secure, SameSite: http.SameSiteLaxMode, MaxAge: -1,
	})
}

// current returns the logged-in username and raw cookie value, if any.
func (s *Sessions) current(r *http.Request) (username, cookieValue string, ok bool) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return "", "", false
	}
	u, ok := s.Verify(c.Value)
	return u, c.Value, ok
}

// --- flash messages -------------------------------------------------------

// setFlash queues a one-time notice ("ok" or "error") shown on the next page.
func (s *Sessions) setFlash(w http.ResponseWriter, kind, msg string) {
	http.SetCookie(w, &http.Cookie{
		Name: flashCookie, Value: base64.RawURLEncoding.EncodeToString([]byte(kind + "|" + msg)),
		Path: "/admin", HttpOnly: true, Secure: s.secure, SameSite: http.SameSiteLaxMode, MaxAge: 60,
	})
}

// takeFlash reads and clears the flash cookie.
func (s *Sessions) takeFlash(w http.ResponseWriter, r *http.Request) (kind, msg string) {
	c, err := r.Cookie(flashCookie)
	if err != nil {
		return "", ""
	}
	http.SetCookie(w, &http.Cookie{Name: flashCookie, Value: "", Path: "/admin", MaxAge: -1})
	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return "", ""
	}
	kind, msg, _ = strings.Cut(string(raw), "|")
	return kind, msg
}
