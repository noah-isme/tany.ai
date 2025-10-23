package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tanydotai/tanyai/backend/internal/db"
)

func TestOpenConnectsToDatabase(t *testing.T) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		t.Skip("POSTGRES_URL not set; skipping database test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := db.Open(ctx, dsn, 3, 3, time.Minute)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	var result int
	err = conn.GetContext(ctx, &result, "SELECT 1")
	require.NoError(t, err)
	require.Equal(t, 1, result)
}
