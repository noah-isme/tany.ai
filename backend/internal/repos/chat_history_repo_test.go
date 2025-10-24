package repos

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/models"
)

func TestChatHistoryRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewChatHistoryRepository(sqlx.NewDb(db, "sqlmock"))

	history := models.ChatHistory{
		ChatID:       uuid.New(),
		UserInput:    "Halo?",
		Model:        "mock-model",
		Prompt:       "prompt",
		PromptHash:   "hash",
		PromptLength: 10,
		ResponseText: "Jawaban",
		LatencyMS:    123,
	}

	rows := sqlmock.NewRows([]string{"id", "chat_id", "user_input", "model", "prompt", "prompt_hash", "prompt_length", "response_text", "latency_ms", "created_at"}).
		AddRow(uuid.New(), history.ChatID, history.UserInput, history.Model, history.Prompt, history.PromptHash, history.PromptLength, history.ResponseText, history.LatencyMS, time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chat_history (chat_id, user_input, model, prompt, prompt_hash, prompt_length, response_text, latency_ms)`)).
		WithArgs(history.ChatID, history.UserInput, history.Model, history.Prompt, history.PromptHash, history.PromptLength, history.ResponseText, history.LatencyMS).
		WillReturnRows(rows)

	if _, err := repo.Create(context.Background(), history); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
