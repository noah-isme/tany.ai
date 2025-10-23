package auth

import "testing"

func TestHashAndComparePassword(t *testing.T) {
	hashed, err := HashPassword("Secret#123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hashed == "Secret#123" {
		t.Fatalf("hashed password should differ from plaintext")
	}
	if err := ComparePassword(hashed, "Secret#123"); err != nil {
		t.Fatalf("compare password: %v", err)
	}
	if err := ComparePassword(hashed, "wrong"); err == nil || err != ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
