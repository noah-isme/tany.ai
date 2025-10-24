package models

import (
	"time"

	"github.com/google/uuid"
)

// ChatHistory represents a persisted chat interaction.
type ChatHistory struct {
	ID           uuid.UUID `db:"id"`
	ChatID       uuid.UUID `db:"chat_id"`
	UserInput    string    `db:"user_input"`
	Model        string    `db:"model"`
	Prompt       string    `db:"prompt"`
	PromptHash   string    `db:"prompt_hash"`
	PromptLength int       `db:"prompt_length"`
	ResponseText string    `db:"response_text"`
	LatencyMS    int       `db:"latency_ms"`
	CreatedAt    time.Time `db:"created_at"`
}
