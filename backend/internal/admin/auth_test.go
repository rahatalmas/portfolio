package admin

import (
	"strings"
	"testing"
)

func TestSessionRoundTrip(t *testing.T) {
	s := NewSessions([]byte("0123456789abcdef0123456789abcdef"), false)
	v := s.Issue("admin")
	u, ok := s.Verify(v)
	if !ok || u != "admin" {
		t.Fatalf("verify failed: %q %v", u, ok)
	}
	// Tampered payload must fail.
	enc, sig, _ := strings.Cut(v, ".")
	if _, ok := s.Verify(enc + "x." + sig); ok {
		t.Fatal("tampered payload accepted")
	}
	// Different secret must fail.
	other := NewSessions([]byte("ffffffffffffffffffffffffffffffff"), false)
	if _, ok := other.Verify(v); ok {
		t.Fatal("cookie accepted under a different secret")
	}
	// CSRF is stable per cookie and differs across cookies.
	if s.CSRF(v) != s.CSRF(v) || s.CSRF(v) == s.CSRF(s.Issue("bob")) {
		t.Fatal("csrf derivation inconsistent")
	}
}

func TestVerifyRejectsGarbage(t *testing.T) {
	s := NewSessions([]byte("0123456789abcdef0123456789abcdef"), false)
	for _, v := range []string{"", "nodot", "a.b", "!!!.!!!"} {
		if _, ok := s.Verify(v); ok {
			t.Fatalf("accepted %q", v)
		}
	}
}
