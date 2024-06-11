package main

import (
	"log"
	"net/http"

	"github.com/flohansen/dasher-server/internal/datastore"
	"github.com/flohansen/dasher-server/internal/routes"
)

func run() error {
	log.Println("listening on port 3000")
	featureStore := datastore.NewInMem()
	return http.ListenAndServe(":3000", routes.New(featureStore))
}

func main() {
	log.Fatal(run())
}
