package dao

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunDatabaseMigrations initializes and runs database migrations
func RunDatabaseMigrations() error {
	m, err := migrate.New(
		"file://dao/migrations",
		"postgres://arran:arran@localhost:5432/arran?sslmode=disable",
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
