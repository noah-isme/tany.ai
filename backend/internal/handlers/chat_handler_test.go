package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/knowledge"
)

func TestHandleChatReturnsMockedAnswer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewChatHandler(knowledge.LoadStaticKnowledgeBase())
	engine := gin.New()
	engine.POST("/chat", handler.HandleChat)

	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBufferString(`{"question":"Layanan apa saja?"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	engine.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
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
}

func TestHandleKnowledgeBaseReturnsSeedData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewChatHandler(knowledge.LoadStaticKnowledgeBase())
	engine := gin.New()
	engine.GET("/kb", handler.HandleKnowledgeBase)

	req := httptest.NewRequest(http.MethodGet, "/kb", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}

	var base knowledge.KnowledgeBase
	if err := json.Unmarshal(res.Body.Bytes(), &base); err != nil {
		t.Fatalf("failed to decode knowledge base: %v", err)
	}

	if base.Profile.Name == "" {
		t.Fatalf("knowledge base should include profile data")
	}
}
