package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	refreshTokenBytes = 32
	clockSkew         = time.Minute
)

// ClockSkew exposes the allowed clock skew for token validation.
const ClockSkew = clockSkew

// ErrInvalidToken indicates a malformed or unverifiable token.
var ErrInvalidToken = errors.New("invalid token")

// ErrTokenExpired indicates the token is no longer valid due to expiration.
var ErrTokenExpired = errors.New("token expired")

// Subject represents identity data embedded within tokens.
type Subject struct {
	ID    uuid.UUID
	Email string
	Roles []string
}

// Claims represents validated access token data.
type Claims struct {
	UserID    uuid.UUID
	Email     string
	Roles     []string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenService issues and validates access/refresh tokens.
type TokenService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
	clock      func() time.Time
}

type accessTokenClaims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

// NewTokenService constructs a TokenService.
func NewTokenService(secret string, accessTTL, refreshTTL time.Duration) (*TokenService, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret must not be empty")
	}
	if accessTTL <= 0 {
		return nil, fmt.Errorf("access token TTL must be positive")
	}
	if refreshTTL <= 0 {
		return nil, fmt.Errorf("refresh token TTL must be positive")
	}
	return &TokenService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		clock:      time.Now,
	}, nil
}

// SetClock overrides the default clock. Intended for testing.
func (s *TokenService) SetClock(clock func() time.Time) {
	if clock == nil {
		s.clock = time.Now
		return
	}
	s.clock = clock
}

// GenerateAccessToken creates a signed JWT for the provided subject.
func (s *TokenService) GenerateAccessToken(sub Subject) (string, error) {
	if sub.ID == uuid.Nil {
		return "", fmt.Errorf("subject ID is required")
	}
	now := s.now()
	claims := accessTokenClaims{
		Email: sub.Email,
		Roles: append([]string(nil), sub.Roles...),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ValidateAccessToken verifies the supplied token and returns claims on success.
func (s *TokenService) ValidateAccessToken(token string) (*Claims, error) {
	parsed := &accessTokenClaims{}
	_, err := jwt.ParseWithClaims(token, parsed, func(_ *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}, jwt.WithTimeFunc(s.now), jwt.WithLeeway(clockSkew))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	subjectStr := parsed.Subject
	if subjectStr == "" {
		return nil, ErrInvalidToken
	}
	subjectID, err := uuid.Parse(subjectStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	roles := append([]string(nil), parsed.Roles...)

	issuedAt := time.Time{}
	if parsed.IssuedAt != nil {
		issuedAt = parsed.IssuedAt.Time
	}
	expiresAt := time.Time{}
	if parsed.ExpiresAt != nil {
		expiresAt = parsed.ExpiresAt.Time
	}

	return &Claims{
		UserID:    subjectID,
		Email:     parsed.Email,
		Roles:     roles,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateRefreshToken creates a secure random token string and hashed value.
func (s *TokenService) GenerateRefreshToken() (token string, hash string, expiresAt time.Time, err error) {
	raw := make([]byte, refreshTokenBytes)
	if _, err = rand.Read(raw); err != nil {
		return "", "", time.Time{}, err
	}
	token = base64.RawURLEncoding.EncodeToString(raw)
	hash = HashRefreshToken(token)
	expiresAt = s.now().Add(s.refreshTTL)
	return
}

// HashRefreshToken hashes the raw refresh token value for storage.
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (s *TokenService) now() time.Time {
	if s.clock != nil {
		return s.clock()
	}
	return time.Now()
}
