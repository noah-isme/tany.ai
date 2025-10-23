package db

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes SQL migration files located in dir against the
// configured PostgreSQL database. The dir can be either an absolute path or a
// path relative to the backend module root.
func RunMigrations(dir, databaseURL string) error {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	source := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.New(source, databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
