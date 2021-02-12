package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/service2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from go server")
	})
	http.ListenAndServe(":8080", nil)
}
