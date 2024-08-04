package main

import (
	"database/sql"
	"flag"
	"log"
	"path"

	"github.com/flohansen/dasher/internal/api"
	"github.com/flohansen/dasher/internal/repository"
	"github.com/flohansen/dasher/internal/routes"
	"github.com/flohansen/dasher/internal/server/feature"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
)

var (
	dataPath = flag.String("data", "/data", "The full path to the sqlite3 database file")
)

func run() error {
	flag.Parse()

	db, err := sql.Open("sqlite3", path.Join(*dataPath, "dasher.db"))
	if err != nil {
		return errors.Wrap(err, "sql open")
	}

	grpcServer := grpc.NewServer()
	store := repository.NewFeatureSQLite(db)
	notifier := feature.NewService(grpcServer, store)
	migrator := repository.NewSQLMigrator(db)
	routes := routes.New(store, notifier)

	return api.New(
		api.WithLogging(),
		api.WithMigrator(migrator),
		api.WithHttpHandler(":3000", routes),
		api.WithNetListenerServer(":50051", grpcServer),
	).Start()
}

func main() {
	log.Fatal(run())
}
