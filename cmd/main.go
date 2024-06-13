package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/flohansen/dasher-server/internal/datastore"
	"github.com/flohansen/dasher-server/internal/routes"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func run() error {
	dataPath := flag.String("data", "/data", "The full path to the sqlite3 database file")
	flag.Parse()

	log.Println("listening on port 3000")
	db, err := sql.Open("sqlite3", path.Join(*dataPath, "dasher.db"))
	if err != nil {
		return errors.Wrap(err, "sql open")
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return errors.Wrap(err, "migrate sqlite3 driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path.Join(*dataPath, "migrations")),
		"sqlite3",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "migrate new")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "migrate up")
	}

	featureStore := datastore.NewSQLite(db)
	return http.ListenAndServe(":3000", routes.New(featureStore))
}

func main() {
	log.Fatal(run())
}
