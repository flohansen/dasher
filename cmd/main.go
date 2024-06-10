package main

import (
	"log"
	"net/http"

	"github.com/flohansen/dasher-server/internal/routes"
)

func run() error {
	return http.ListenAndServe(":3000", routes.New())
}

func main() {
	log.Fatal(run())
}
