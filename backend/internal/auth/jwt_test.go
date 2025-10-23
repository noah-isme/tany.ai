package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenServiceAccessTokenLifecycle(t *testing.T) {
	svc, err := NewTokenService(strings.Repeat("s", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("new token service: %v", err)
	}
	subject := Subject{ID: uuid.New(), Email: "admin@example.com", Roles: []string{"admin"}}
	token, err := svc.GenerateAccessToken(subject)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}
	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("validate access token: %v", err)
	}
	if claims.UserID != subject.ID {
		t.Fatalf("expected subject ID %s, got %s", subject.ID, claims.UserID)
	}
	if claims.Email != subject.Email {
		t.Fatalf("expected email %s, got %s", subject.Email, claims.Email)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
		t.Fatalf("expected roles [admin], got %v", claims.Roles)
	}
}

func TestTokenServiceExpiredAccessToken(t *testing.T) {
	svc, err := NewTokenService(strings.Repeat("t", 64), time.Second, time.Hour)
	if err != nil {
		t.Fatalf("new token service: %v", err)
	}
	base := time.Now()
	svc.SetClock(func() time.Time { return base })
	token, err := svc.GenerateAccessToken(Subject{ID: uuid.New(), Email: "user@example.com"})
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}
	svc.SetClock(func() time.Time { return base.Add(clockSkew + 2*time.Second) })
	if _, err := svc.ValidateAccessToken(token); err != ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	svc, err := NewTokenService(strings.Repeat("r", 64), time.Minute, time.Hour)
	if err != nil {
		t.Fatalf("new token service: %v", err)
	}
	token, hash, expiresAt, err := svc.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}
	if token == "" {
		t.Fatal("expected refresh token value")
	}
	if hash != HashRefreshToken(token) {
		t.Fatalf("hash mismatch: got %s", hash)
	}
	if !expiresAt.After(time.Now().Add(30 * time.Minute)) {
		t.Fatal("expected refresh token expiry in the future")
	}
}
