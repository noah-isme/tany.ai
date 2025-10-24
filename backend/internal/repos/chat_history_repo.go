package repos

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

// ChatHistoryRepository persists chat interactions for auditing and analytics.
type ChatHistoryRepository interface {
	Create(ctx context.Context, history models.ChatHistory) (models.ChatHistory, error)
	ListRecentByChat(ctx context.Context, chatID uuid.UUID, limit int) ([]models.ChatHistory, error)
}

// NewChatHistoryRepository constructs a SQL-backed repository.
func NewChatHistoryRepository(db *sqlx.DB) ChatHistoryRepository {
	return &chatHistoryRepository{db: db}
}

type chatHistoryRepository struct {
	db *sqlx.DB
}

func (r *chatHistoryRepository) Create(ctx context.Context, history models.ChatHistory) (models.ChatHistory, error) {
	const query = `INSERT INTO chat_history (chat_id, user_input, model, prompt, prompt_hash, prompt_length, response_text, latency_ms)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, chat_id, user_input, model, prompt, prompt_hash, prompt_length, response_text, latency_ms, created_at`

	var created models.ChatHistory
	if err := r.db.GetContext(ctx, &created, query,
		history.ChatID,
		history.UserInput,
		history.Model,
		history.Prompt,
		history.PromptHash,
		history.PromptLength,
		history.ResponseText,
		history.LatencyMS,
	); err != nil {
		return models.ChatHistory{}, err
	}
	return created, nil
}

func (r *chatHistoryRepository) ListRecentByChat(ctx context.Context, chatID uuid.UUID, limit int) ([]models.ChatHistory, error) {
	if limit <= 0 {
		limit = 5
	}
	const query = `SELECT id, chat_id, user_input, model, prompt, prompt_hash, prompt_length, response_text, latency_ms, created_at
FROM chat_history WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2`
	var rows []models.ChatHistory
	if err := r.db.SelectContext(ctx, &rows, query, chatID, limit); err != nil {
		return nil, err
	}
	return rows, nil
}
