package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/flohansen/dasher-server/internal/sqlc"
)

type GetFeatureResponse struct {
	FeatureID   string `json:"featureId"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (routes *Routes) getFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := routes.store.GetAll(context.Background())
	if err != nil {
		log.Printf("error while fetching all features: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]GetFeatureResponse, len(features))
	for i, feature := range features {
		response[i] = GetFeatureResponse(feature)
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error while encoding features: %v", err)
	}
}

type PostFeatureRequest struct {
	FeatureID   string `json:"featureId"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (routes *Routes) postFeature(w http.ResponseWriter, r *http.Request) {
	var req PostFeatureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error while decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := routes.store.Upsert(context.Background(), sqlc.Feature(req)); err != nil {
		log.Printf("error while upserting feature: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	routes.notifier.Notify(sqlc.Feature(req))
}

func (routes *Routes) deleteFeature(w http.ResponseWriter, r *http.Request) {
	featureID := r.PathValue("featureId")
	if err := routes.store.Delete(context.Background(), featureID); err != nil {
		log.Printf("error while deleting feature: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
