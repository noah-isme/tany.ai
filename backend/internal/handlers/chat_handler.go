package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/knowledge"
)

// ChatHandler exposes HTTP handlers for chat related endpoints.
type ChatHandler struct {
	base knowledge.KnowledgeBase
}

// ChatRequest represents the incoming chat payload.
type ChatRequest struct {
	Question string `json:"question" binding:"required"`
}

// ChatResponse is a simplified response structure used during the MVP stage.
type ChatResponse struct {
	Answer string                  `json:"answer"`
	Prompt string                  `json:"prompt"`
	Scope  knowledge.KnowledgeBase `json:"scope"`
}

// NewChatHandler returns a handler configured with the provided knowledge base.
func NewChatHandler(base knowledge.KnowledgeBase) *ChatHandler {
	return &ChatHandler{base: base}
}

// HandleChat processes the chat question and returns a mock answer grounded by the knowledge base.
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var payload ChatRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "question field is required",
		})
		return
	}

	response := ChatResponse{
		Answer: knowledge.SummarizeForHuman(payload.Question, h.base),
		Prompt: knowledge.BuildSystemPrompt(h.base),
		Scope:  h.base,
	}

	c.JSON(http.StatusOK, response)
}

// HandleKnowledgeBase exposes the static knowledge base for debugging/testing.
func (h *ChatHandler) HandleKnowledgeBase(c *gin.Context) {
	c.JSON(http.StatusOK, h.base)
}
