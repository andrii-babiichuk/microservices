package handlers

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/kpi-lab/microservices-demo/services/service1/repository"
)

type Server struct {
	db repository.Visits
}

func New(db repository.Visits) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	var current int
	var err error
	log.Println("ping request")

	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		}
	}()

	err = s.db.Inc(r.Context())
	if err != nil {
		log.Println("failed to increment visits: %w", err)
	}
	current, err = s.db.Get(r.Context())
	if err != nil {
		log.Println("failed to get visits: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("Service1 is healthy, visited: %d times", current)
	_, err = w.Write([]byte(msg))
}
