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

func (repo *SQLiteDatastore) Upsert(ctx context.Context, feature sqlc.Feature) error {
	return repo.q.UpsertFeature(ctx, sqlc.UpsertFeatureParams(feature))
}

func (repo *SQLiteDatastore) Delete(ctx context.Context, featureID string) error {
	return repo.q.DeleteFeature(ctx, featureID)
}
