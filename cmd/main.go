package main

import (
	"log"
	"net/http"

	"github.com/flohansen/dasher-server/internal/routes"
)

func run() error {
	log.Println("listening on port 3000")
	return http.ListenAndServe(":3000", routes.New())
}

func main() {
	log.Fatal(run())
}
