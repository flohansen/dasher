package datastore

import (
	"context"

	"github.com/flohansen/dasher-server/internal/sqlc"
)

type SQLiteDatastore struct {
	q *sqlc.Queries
}

func NewSQLite(db sqlc.DBTX) *SQLiteDatastore {
	q := sqlc.New(db)
	return &SQLiteDatastore{q}
}

func (repo *SQLiteDatastore) GetAll(ctx context.Context) ([]sqlc.Feature, error) {
	return repo.q.GetAllFeatures(ctx)
}
