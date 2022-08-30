package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/kpi-lab/ci/services/service2/handlers"
)

const httpPort = 8080

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/service2/ping", handlers.Ping)

	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r)
	if err != nil {
		log.Fatal(err)
	}
}
