package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

type stubKnowledge struct {
	base kb.KnowledgeBase
}

func (s *stubKnowledge) Get(ctx context.Context) (kb.KnowledgeBase, string, bool, error) {
	return s.base, "\"etag\"", false, nil
}

func (s *stubKnowledge) CacheTTL() time.Duration {
	return time.Minute
}

type historyRecorder struct {
	records []models.ChatHistory
}

func (h *historyRecorder) Create(ctx context.Context, history models.ChatHistory) (models.ChatHistory, error) {
	history.ID = models.ChatHistory{}.ID
	history.CreatedAt = time.Now()
	h.records = append(h.records, history)
	return history, nil
}

func (h *historyRecorder) ListRecentByChat(ctx context.Context, chatID uuid.UUID, limit int) ([]models.ChatHistory, error) {
	return nil, nil
}

var _ repos.ChatHistoryRepository = (*historyRecorder)(nil)

func TestHandleChatStoresHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	knowledge := &stubKnowledge{base: kb.KnowledgeBase{
		Profile:  kb.Profile{Name: "Tanya"},
		Services: []kb.Service{{Name: "Consulting"}},
		Projects: []kb.Project{{Title: "Featured", IsFeatured: true}},
	}}
	history := &historyRecorder{}
	handler := NewChatHandler(knowledge, history, "mock-model")
	engine := gin.New()
	engine.POST("/chat", handler.HandleChat)

	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBufferString(`{"question":"Layanan apa saja?"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	if len(history.records) != 1 {
		t.Fatalf("expected history to be stored")
	}

	var payload ChatResponse
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Answer == "" {
		t.Fatalf("expected answer to be populated")
	}
	if payload.Prompt == "" {
		t.Fatalf("prompt should not be empty")
	}
	if payload.Model != "mock-model" {
		t.Fatalf("expected model to be propagated")
	}
}

func TestHandleKnowledgeBaseSetsCachingHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	knowledge := &stubKnowledge{base: kb.KnowledgeBase{Profile: kb.Profile{Name: "Tanya"}}}
	handler := NewChatHandler(knowledge, nil, "mock")
	engine := gin.New()
	engine.GET("/kb", handler.HandleKnowledgeBase)

	req := httptest.NewRequest(http.MethodGet, "/kb", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	if etag := res.Header().Get("ETag"); etag == "" {
		t.Fatalf("expected ETag header to be set")
	}
	if cacheControl := res.Header().Get("Cache-Control"); cacheControl == "" {
		t.Fatalf("expected Cache-Control header to be set")
	}
}
