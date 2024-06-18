package datastore

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/pkg/errors"
)

type SQLMigrator struct {
	db *sql.DB
}

func NewSQLMigrator(db *sql.DB) *SQLMigrator {
	return &SQLMigrator{db}
}

func (mig *SQLMigrator) Migrate() error {
	driver, err := sqlite3.WithInstance(mig.db, &sqlite3.Config{})
	if err != nil {
		return errors.Wrap(err, "migrate sqlite3 driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
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
