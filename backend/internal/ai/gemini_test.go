package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeminiGenerateSuccess(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request")
		}
		if got := r.URL.Query().Get("key"); got != "secret" {
			t.Fatalf("expected key query param, got %q", got)
		}
		response := map[string]any{
			"candidates": []any{
				map[string]any{
					"content": map[string]any{
						"parts": []any{map[string]any{"text": "Hello"}},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	gemini := NewGemini("secret", "custom-model")
	gemini.Endpoint = server.URL
	gemini.Client = server.Client()

	resp, err := gemini.Generate(context.Background(), Request{Prompt: "Hi", MaxTokens: 16})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Text != "Hello" {
		t.Fatalf("expected response text, got %q", resp.Text)
	}
}

func TestGeminiGenerateHandlesEmptyResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"candidates": []any{}})
	})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	gemini := NewGemini("secret", "")
	gemini.Endpoint = server.URL
	gemini.Client = server.Client()

	if _, err := gemini.Generate(context.Background(), Request{Prompt: "Hi"}); err == nil {
		t.Fatalf("expected error for empty response")
	}
}

func TestGeminiGenerateRequiresAPIKey(t *testing.T) {
	gemini := NewGemini("", "")
	if _, err := gemini.Generate(context.Background(), Request{Prompt: "Hi"}); err == nil {
		t.Fatalf("expected missing key error")
	}
}
