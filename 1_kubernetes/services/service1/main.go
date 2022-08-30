package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/kpi-lab/ci/services/service1/handlers"
)

const httpPort = 8080

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/service1/ping", handlers.Ping)

	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
