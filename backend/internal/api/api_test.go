package api

import (
	"testing"
	"time"
)

func TestValidateMessage(t *testing.T) {
	ok := messageInput{Name: "Ada", Email: "ada@example.com", Message: "Hello"}
	if p := validateMessage(ok); p != "" {
		t.Fatalf("valid input rejected: %s", p)
	}
	bad := []messageInput{
		{Email: "a@b.co", Message: "x"},
		{Name: "a", Message: "x"},
		{Name: "a", Email: "a@b.co"},
		{Name: "a", Email: "not-an-email", Message: "x"},
	}
	for i, in := range bad {
		if p := validateMessage(in); p == "" {
			t.Fatalf("case %d should fail", i)
		}
	}
}

func TestRateLimiter(t *testing.T) {
	l := newRateLimiter(2, 50*time.Millisecond)
	if !l.Allow("a") || !l.Allow("a") {
		t.Fatal("first two should pass")
	}
	if l.Allow("a") {
		t.Fatal("third should be limited")
	}
	if !l.Allow("b") {
		t.Fatal("other key unaffected")
	}
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("a") {
		t.Fatal("window should have reset")
	}
}
