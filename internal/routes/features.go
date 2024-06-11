package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func (routes *Routes) getFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := routes.featureStore.GetAll()
	if err != nil {
		log.Printf("error while fetching all features: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(features); err != nil {
		log.Printf("error while encoding features: %v", err)
	}
}
