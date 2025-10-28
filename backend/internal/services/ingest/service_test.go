package ingest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestServiceSyncProjects(t *testing.T) {
	etagValue := `"abc123"`
	projectHTML := `<!DOCTYPE html><html><head><script type="application/ld+json">{"@context":"https://schema.org","@type":"CreativeWork","name":"Test Project","headline":"A summary","about":"Detailed body","datePublished":"2024-02-10","image":"https://example.com/image.png"}</script></head><body></body></html>`

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/robots.txt":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("User-agent: *\nAllow: /\n"))
		case "/sitemap.xml":
			if inm := r.Header.Get("If-None-Match"); inm == etagValue {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set("ETag", etagValue)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>` + server.URL + `/project/test/</loc></url></urlset>`))
		case "/project/test/":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(projectHTML))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	baseURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}

	service := NewService(2*time.Second, 120, []string{baseURL.Host})
	source := Source{ID: uuid.New(), Name: "Test", BaseURL: baseURL, SourceType: "auto"}

	result, err := service.Sync(context.Background(), source)
	if err != nil {
		t.Fatalf("sync error: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	item := result.Items[0]
	if item.Kind != "project" {
		t.Fatalf("unexpected kind %s", item.Kind)
	}
	if item.Title != "Test Project" {
		t.Fatalf("unexpected title %s", item.Title)
	}
	if item.Summary == nil || *item.Summary != "A summary" {
		t.Fatalf("unexpected summary %+v", item.Summary)
	}
	if item.Content == nil || *item.Content != "Detailed body" {
		t.Fatalf("unexpected content %+v", item.Content)
	}
	if item.Metadata == nil || item.Metadata["image"] != "https://example.com/image.png" {
		t.Fatalf("expected image metadata, got %+v", item.Metadata)
	}
	if result.ETag == nil || *result.ETag != etagValue {
		t.Fatalf("expected etag in result")
	}

	// Second sync should short-circuit via 304
	source.ETag = result.ETag
	if _, err := service.Sync(context.Background(), source); err != ErrNotModified {
		t.Fatalf("expected ErrNotModified, got %v", err)
	}
}
