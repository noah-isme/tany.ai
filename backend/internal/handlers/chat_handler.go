package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tanydotai/tanyai/backend/internal/httpapi"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
	"github.com/tanydotai/tanyai/backend/internal/services/prompt"
)

// KnowledgeService defines the behaviour required from a knowledge base provider.
type KnowledgeService interface {
	Get(ctx context.Context) (kb.KnowledgeBase, string, bool, error)
	CacheTTL() time.Duration
}

// ChatHandler exposes HTTP handlers for chat and knowledge base endpoints.
type ChatHandler struct {
	knowledge KnowledgeService
	history   repos.ChatHistoryRepository
	modelName string
}

// ChatRequest represents the incoming chat payload.
type ChatRequest struct {
	Question string `json:"question" binding:"required"`
	ChatID   string `json:"chatId"`
}

// ChatResponse contains the assistant reply and metadata returned to clients.
type ChatResponse struct {
	ChatID string `json:"chatId"`
	Answer string `json:"answer"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// NewChatHandler constructs a ChatHandler with the provided dependencies.
func NewChatHandler(knowledge KnowledgeService, history repos.ChatHistoryRepository, modelName string) *ChatHandler {
	return &ChatHandler{knowledge: knowledge, history: history, modelName: modelName}
}

// HandleChat processes the chat question and stores the interaction history.
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var payload ChatRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "question field is required", nil)
		return
	}

	base, _, cacheHit, err := h.knowledge.Get(c.Request.Context())
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load knowledge base", nil)
		return
	}
	c.Set("kb_cache_hit", cacheHit)
	c.Set("model", h.modelName)

	promptText := prompt.BuildSystemPrompt(base)
	answer := prompt.SummarizeForHuman(payload.Question, base)
	promptHash := sha256.Sum256([]byte(promptText))
	promptLength := len([]rune(promptText))

	chatID := uuid.New()
	if payload.ChatID != "" {
		parsed, err := uuid.Parse(payload.ChatID)
		if err != nil {
			httpapi.RespondError(c, http.StatusBadRequest, httpapi.ErrorCodeValidation, "chatId must be a valid UUID", nil)
			return
		}
		chatID = parsed
	}
	c.Set("chat_id", chatID.String())
	c.Set("prompt_length", promptLength)

	record := models.ChatHistory{
		ChatID:       chatID,
		UserInput:    payload.Question,
		Model:        h.modelName,
		Prompt:       promptText,
		PromptHash:   hex.EncodeToString(promptHash[:]),
		PromptLength: promptLength,
		ResponseText: answer,
		LatencyMS:    0,
		CreatedAt:    time.Now(),
	}

	if h.history != nil {
		if _, err := h.history.Create(c.Request.Context(), record); err != nil {
			httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to store chat history", nil)
			return
		}
	}

	response := ChatResponse{
		ChatID: chatID.String(),
		Answer: answer,
		Model:  h.modelName,
		Prompt: promptText,
	}

	c.JSON(http.StatusOK, response)
}

// HandleKnowledgeBase exposes the aggregated knowledge base with caching headers.
func (h *ChatHandler) HandleKnowledgeBase(c *gin.Context) {
	data, etag, cacheHit, err := h.knowledge.Get(c.Request.Context())
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, httpapi.ErrorCodeInternal, "failed to load knowledge base", nil)
		return
	}

	if match := c.GetHeader("If-None-Match"); match != "" && match == etag {
		c.Header("ETag", etag)
		c.Status(http.StatusNotModified)
		return
	}

	maxAge := int(h.knowledge.CacheTTL().Seconds())
	if maxAge <= 0 {
		maxAge = 60
	}
	c.Header("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
	c.Header("ETag", etag)
	if cacheHit {
		c.Header("X-Cache", "HIT")
	} else {
		c.Header("X-Cache", "MISS")
	}
	c.Set("kb_cache_hit", cacheHit)

	c.JSON(http.StatusOK, data)
}
