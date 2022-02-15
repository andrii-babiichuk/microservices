package handlers

import (
	"log"
	"net/http"
)

func Ping(w http.ResponseWriter, _ *http.Request) {
	log.Println("ping request")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Service2 is healthy"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
