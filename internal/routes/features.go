package routes

import (
	"context"
	"database/sql"
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
	features, err := routes.featureStore.GetAll(context.Background())
	if err != nil {
		log.Printf("error while fetching all features: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]GetFeatureResponse, len(features))
	for i, feature := range features {
		response[i] = GetFeatureResponse{
			FeatureID:   feature.FeatureID,
			Description: feature.Description.String,
			Enabled:     feature.Enabled.Bool,
		}
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

	if err := routes.featureStore.Upsert(context.Background(), sqlc.Feature{
		FeatureID:   req.FeatureID,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Enabled:     sql.NullBool{Bool: req.Enabled, Valid: req.Enabled},
	}); err != nil {
		log.Printf("error while upserting feature: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (routes *Routes) deleteFeature(w http.ResponseWriter, r *http.Request) {
	featureID := r.PathValue("featureId")
	if err := routes.featureStore.Delete(context.Background(), featureID); err != nil {
		log.Printf("error while deleting feature: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
