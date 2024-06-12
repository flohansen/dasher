package routes

import (
	"context"
	"net/http"

	"github.com/flohansen/dasher-server/internal/sqlc"
)

type FeatureStore interface {
	GetAll(ctx context.Context) ([]sqlc.Feature, error)
}

type Routes struct {
	mux          *http.ServeMux
	featureStore FeatureStore
}

func New(featureStore FeatureStore) *Routes {
	routes := Routes{
		mux:          http.NewServeMux(),
		featureStore: featureStore,
	}

	routes.mux.HandleFunc("GET /api/v1/features", routes.getFeatures)
	return &routes
}

func (routes *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.mux.ServeHTTP(w, r)
}
