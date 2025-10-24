package kb

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestAggregatorLoadsDataAndCaches(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	columnsProfile := []string{"id", "name", "title", "bio", "email", "phone", "location", "avatar_url", "updated_at"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, title, bio, email, phone, location, avatar_url, updated_at FROM profile ORDER BY updated_at DESC LIMIT 1`)).
		WillReturnRows(sqlmock.NewRows(columnsProfile).AddRow("00000000-0000-0000-0000-000000000001", "Tanya", "Lead", "Bio", "hello@tany.ai", "", "Jakarta", "", time.Now()))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT name FROM skills ORDER BY "order" ASC, name ASC`)).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Go"))

	columnsServices := []string{"id", "name", "description", "price_min", "price_max", "currency", "duration_label", "order"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, description, price_min, price_max, currency, duration_label, "order" FROM services WHERE is_active = TRUE ORDER BY "order" ASC, name ASC`)).
		WillReturnRows(sqlmock.NewRows(columnsServices).AddRow("00000000-0000-0000-0000-000000000010", "Dev", "Desc", 1000.0, 2000.0, "IDR", "2 minggu", 1))

	columnsProjects := []string{"id", "title", "description", "tech_stack", "project_url", "category", "order", "is_featured"}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, tech_stack, project_url, category, "order", is_featured FROM projects ORDER BY is_featured DESC, "order" ASC, title ASC`)).
		WillReturnRows(sqlmock.NewRows(columnsProjects).AddRow("00000000-0000-0000-0000-000000000020", "Proj", "Impact", `{"Go"}`, "https://example.com", "Web", 1, true))

	aggregator := NewAggregator(sqlx.NewDb(db, "sqlmock"), time.Second)

	data, etag, hit, err := aggregator.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hit {
		t.Fatalf("expected first call to be a cache miss")
	}
	if etag == "" {
		t.Fatalf("etag should not be empty")
	}
	if data.Profile.Name != "Tanya" {
		t.Fatalf("expected profile to be loaded")
	}

	// Second call should reuse cache, thus no additional expectations.
	data2, _, hit, err := aggregator.Get(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on second get: %v", err)
	}
	if !hit {
		t.Fatalf("expected cache hit on second call")
	}
	if data2.Profile.Name != data.Profile.Name {
		t.Fatalf("expected cached data to match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
