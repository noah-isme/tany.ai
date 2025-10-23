package seed_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/tanydotai/tanyai/backend/internal/db"
	"github.com/tanydotai/tanyai/backend/internal/seed"
)

func TestSeedInsertsSampleData(t *testing.T) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		t.Skip("POSTGRES_URL not set; skipping seeder test")
	}

	migrationsDir := filepath.Join("..", "..", "migrations")
	require.NoError(t, db.RunMigrations(migrationsDir, dsn))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := db.Open(ctx, dsn, 3, 3, time.Minute)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	require.NoError(t, seed.Seed(ctx, conn))

	assertCountGreaterThanZero(t, conn, "profile")
	assertCountGreaterThanZero(t, conn, "skills")
	assertCountGreaterThanZero(t, conn, "services")
	assertCountGreaterThanZero(t, conn, "projects")

	var featuredCount int
	require.NoError(t, conn.Get(&featuredCount, "SELECT COUNT(*) FROM projects WHERE is_featured = true"))
	require.Greater(t, featuredCount, 0)
}

func assertCountGreaterThanZero(t *testing.T, dbConn *sqlx.DB, table string) {
	t.Helper()
	var count int
	query := "SELECT COUNT(*) FROM " + table
	require.NoError(t, dbConn.Get(&count, query))
	require.Greater(t, count, 0)
}
