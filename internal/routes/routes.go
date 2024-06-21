package routes

import (
	"context"
	"net/http"

	"github.com/flohansen/dasher/internal/sqlc"
)

type FeatureStore interface {
	GetAll(ctx context.Context) ([]sqlc.Feature, error)
	Upsert(ctx context.Context, feature sqlc.Feature) error
	Delete(ctx context.Context, featureID string) error
}

type Notifier interface {
	Notify(sqlc.Feature) error
}

type Routes struct {
	mux      *http.ServeMux
	store    FeatureStore
	notifier Notifier
}

func New(store FeatureStore, notifier Notifier) *Routes {
	routes := Routes{
		mux:      http.NewServeMux(),
		store:    store,
		notifier: notifier,
	}

	routes.mux.HandleFunc("GET /api/v1/features", routes.getFeatures)
	routes.mux.HandleFunc("POST /api/v1/features", routes.postFeature)
	routes.mux.HandleFunc("DELETE /api/v1/features/{featureId}", routes.deleteFeature)
	return &routes
}

func (routes *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.mux.ServeHTTP(w, r)
}
