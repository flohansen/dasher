package repository

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type SQLiteMigrator struct {
	db *sql.DB
}

func NewSQLiteMigrator(db *sql.DB) *SQLiteMigrator {
	return &SQLiteMigrator{db}
}

func (mig *SQLiteMigrator) Migrate() error {
	driver, err := sqlite3.WithInstance(mig.db, &sqlite3.Config{})
	if err != nil {
		return errors.Wrap(err, "migrate sqlite3 driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "migrate new")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "migrate up")
	}

	return nil
}
