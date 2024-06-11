package routes

import "net/http"

type Routes struct {
	mux *http.ServeMux
}

func New() *Routes {
	routes := Routes{
		mux: http.NewServeMux(),
	}

	routes.mux.HandleFunc("GET /api/v1/features", routes.getFeatures)
	return &routes
}

func (routes *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.mux.ServeHTTP(w, r)
}
