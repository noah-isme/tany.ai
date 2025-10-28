package seed

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestSeederUsesEmbeddedData(t *testing.T) {
	t.Setenv("SEED_DATA_PATH", "")
	payload := mustLoadPayload(t, defaultDataFS)

	dbx, mock := setupSeederMock(t, payload)
	seeder := NewSeeder(dbx)

	err := seeder.Seed(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSeederHonorsSeedDataPathOverride(t *testing.T) {
	tmp := t.TempDir()
	copySeedFiles(t, tmp)
	t.Setenv("SEED_DATA_PATH", tmp)

	payload := mustLoadPayload(t, os.DirFS(tmp))

	dbx, mock := setupSeederMock(t, payload)
	seeder := NewSeeder(dbx)

	err := seeder.Seed(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func setupSeederMock(t *testing.T, payload seedPayload) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
	})

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO profile").
		WithArgs(payload.Profile.ID, payload.Profile.Name, payload.Profile.Title, payload.Profile.Bio, payload.Profile.Email, payload.Profile.Phone, payload.Profile.Location, payload.Profile.AvatarURL, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	for _, item := range payload.Skills {
		mock.ExpectExec("INSERT INTO skills").
			WithArgs(item.ID, item.Name, item.Order).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	for _, item := range payload.Services {
		mock.ExpectExec("INSERT INTO services").
			WithArgs(item.ID, item.Name, item.Description, item.PriceMin, item.PriceMax, item.Currency, item.DurationLabel, item.IsActive, item.Order).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	for _, item := range payload.Projects {
		mock.ExpectExec("INSERT INTO projects").
			WithArgs(item.ID, item.Title, item.Description, item.TechStack, item.ImageURL, item.ProjectURL, item.Category, item.DurationLabel, item.PriceLabel, item.BudgetLabel, item.Order, item.IsFeatured).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	for range payload.Leads {
		mock.ExpectExec("INSERT INTO leads").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(payload.Admin.ID, payload.Admin.Email, sqlmock.AnyArg(), payload.Admin.Name).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("INSERT INTO user_roles").
		WithArgs(payload.Admin.ID, payload.Admin.Role).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	return sqlx.NewDb(db, "sqlmock"), mock
}

func mustLoadPayload(t *testing.T, fsys fs.FS) seedPayload {
	t.Helper()
	payload, err := loadSeedData(fsys)
	require.NoError(t, err)
	return payload
}

func copySeedFiles(t *testing.T, dst string) {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "unable to determine caller")
	srcDir := filepath.Join(filepath.Dir(file), "data")

	entries := []string{
		"profile.json",
		"skills.json",
		"services.json",
		"projects.json",
		"leads.json",
		"admin_user.json",
	}

	for _, name := range entries {
		content, err := os.ReadFile(filepath.Join(srcDir, name))
		require.NoErrorf(t, err, "read %s", name)
		err = os.WriteFile(filepath.Join(dst, name), content, 0o644)
		require.NoErrorf(t, err, "write %s", name)
	}
}
