package routes

import "net/http"

type Routes struct {
	mux *http.ServeMux
}

func New() *Routes {
	mux := http.NewServeMux()
	return &Routes{mux}
}

func (routes *Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes.mux.ServeHTTP(w, r)
}
