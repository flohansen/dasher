//go:generate sqlc -f ../../sqlc.yaml generate

package repository

import (
	"context"

	"github.com/flohansen/dasher/internal/sqlc"
)

type FeatureSQLite struct {
	q *sqlc.Queries
}

func NewFeatureSQLite(db sqlc.DBTX) *FeatureSQLite {
	q := sqlc.New(db)
	return &FeatureSQLite{q}
}

func (repo *FeatureSQLite) GetAll(ctx context.Context) ([]sqlc.Feature, error) {
	return repo.q.GetAllFeatures(ctx)
}

func (repo *FeatureSQLite) Upsert(ctx context.Context, feature sqlc.Feature) error {
	return repo.q.UpsertFeature(ctx, sqlc.UpsertFeatureParams(feature))
}

func (repo *FeatureSQLite) Delete(ctx context.Context, featureID string) error {
	return repo.q.DeleteFeature(ctx, featureID)
}
