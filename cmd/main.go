package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/flohansen/dasher-server/internal/datastore"
	"github.com/flohansen/dasher-server/internal/routes"
	"github.com/pkg/errors"

	_ "github.com/mattn/go-sqlite3"
)

func run() error {
	log.Println("listening on port 3000")
	db, err := sql.Open("sqlite3", "./dasher.db")
	if err != nil {
		return errors.Wrap(err, "sql open")
	}

	featureStore := datastore.NewSQLite(db)
	return http.ListenAndServe(":3000", routes.New(featureStore))
}

func main() {
	log.Fatal(run())
}
