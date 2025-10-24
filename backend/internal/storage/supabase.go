package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/tanydotai/tanyai/backend/internal/config"
)

// SupabaseStorage stores files using Supabase Storage buckets.
type SupabaseStorage struct {
	client     *http.Client
	baseURL    string
	bucket     string
	serviceKey string
	publicBase string
}

// NewSupabaseStorage creates a SupabaseStorage instance.
func NewSupabaseStorage(cfg config.SupabaseConfig) (*SupabaseStorage, error) {
	if cfg.URL == "" || cfg.Bucket == "" || cfg.ServiceRole == "" {
		return nil, fmt.Errorf("storage: incomplete supabase configuration")
	}
	client := &http.Client{Timeout: 30 * time.Second}
	return &SupabaseStorage{
		client:     client,
		baseURL:    strings.TrimSuffix(cfg.URL, "/"),
		bucket:     cfg.Bucket,
		serviceKey: cfg.ServiceRole,
		publicBase: strings.TrimSuffix(cfg.PublicURL, "/"),
	}, nil
}

// Put uploads content to Supabase Storage and returns a public URL.
func (s *SupabaseStorage) Put(ctx context.Context, key string, content []byte, contentType string) (string, error) {
	if len(content) == 0 {
		return "", fmt.Errorf("storage: empty content")
	}
	endpoint := s.buildObjectURL(key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(content))
	if err != nil {
		return "", fmt.Errorf("storage: build supabase request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "false")
	req.ContentLength = int64(len(content))

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("storage: supabase upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf("storage: supabase upload error %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return s.publicURL(key), nil
}

func (s *SupabaseStorage) buildObjectURL(key string) string {
	cleanedKey := strings.TrimPrefix(key, "/")
	joined := path.Join("storage/v1/object", s.bucket, cleanedKey)
	return s.baseURL + "/" + joined
}

func (s *SupabaseStorage) publicURL(key string) string {
	base := s.publicBase
	if base == "" {
		base = s.baseURL + "/storage/v1/object/public/" + s.bucket
	}
	return strings.TrimSuffix(base, "/") + "/" + strings.TrimPrefix(key, "/")
}
