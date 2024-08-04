package main

import (
	"database/sql"
	"flag"
	"log"
	"net"
	"net/http"
	"path"

	"github.com/flohansen/dasher/internal/repository"
	"github.com/flohansen/dasher/internal/routes"
	"github.com/flohansen/dasher/internal/server/feature"
	"github.com/flohansen/dasher/pkg/proto"
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

	migrator := repository.NewSQLiteMigrator(db)
	if err := migrator.Migrate(); err != nil {
		return errors.Wrap(err, "migrate")
	}

	store := repository.NewFeatureSQLite(db)
	featureService := feature.NewService(store)
	errs := make(chan error, 1)

	go startGrpcServer(featureService, errs)
	go startHttpServer(store, featureService, errs)
	return <-errs
}

func startGrpcServer(featureService *feature.Service, errs chan error) {
	grpcServer := grpc.NewServer()
	proto.RegisterFeatureStateServiceServer(grpcServer, featureService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		errs <- err
		return
	}

	errs <- grpcServer.Serve(lis)
}

func startHttpServer(store *repository.FeatureSQLite, featureService *feature.Service, errs chan error) {
	routes := routes.New(store, featureService)
	errs <- http.ListenAndServe(":3000", routes)
}

func main() {
	log.Fatalf("fatal error: %s", run())
}
