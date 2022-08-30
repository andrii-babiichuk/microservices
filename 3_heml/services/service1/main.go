package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"

	"gitlab.com/kpi-lab/microservices-demo/services/service1/handlers"
	"gitlab.com/kpi-lab/microservices-demo/services/service1/repository/postgres"
)

var (
	httpPort int
	pgHost   string
	pgUser   string
	pgPass   string
	pgDb     string
)

func init() {
	httpPort = 8080
	pgUser = os.Getenv("POSTGRES_USER")
	pgPass = os.Getenv("POSTGRES_PASSWORD")
	pgHost = os.Getenv("POSTGRES_HOST")
	pgDb = os.Getenv("POSTGRES_DB")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnector := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPass, pgHost, pgDb)

	conn, err := pgx.Connect(ctx, dbConnector)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	visits := postgres.New(conn)
	server := handlers.New(visits)

	r := mux.NewRouter()
	r.HandleFunc("/api/service1/ping", server.Ping)

	err = http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
