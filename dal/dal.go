package dal

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jicol-95/arran/config"
	_ "github.com/lib/pq"
)

func InitDB(cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", buildConnectionString(cfg))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunDatabaseMigrations(cfg config.PostgresConfig) error {
	fmt.Printf("Running migrations on %s\n", buildConnectionString(cfg))
	m, err := migrate.New(
		"file://dal/migration",
		buildConnectionString(cfg),
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

func buildConnectionString(config config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}
