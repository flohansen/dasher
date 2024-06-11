package routes

import "net/http"

func (routes *Routes) getFeatures(w http.ResponseWriter, r *http.Request) {
	_, err := routes.featureStore.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
