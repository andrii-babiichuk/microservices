package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getRemote(url string, results chan string) {
	resp, err := http.Get(url)
	if err != nil {
		results <- err.Error()
	} else {
		defer resp.Body.Close()
		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			results <- err2.Error()
		} else {
			results <- string(body)
		}
	}
}

func main() {
	http.HandleFunc("/api/root-service", func(w http.ResponseWriter, r *http.Request) {
		service1Result := make(chan string)
		service2Result := make(chan string)

		go getRemote("http://service1-service/api/service1", service1Result)
		go getRemote("http://service2-service/api/service2", service2Result)

		response1Message := fmt.Sprintf("service1 response: %s \n\n", <-service1Result)
		response2Message := fmt.Sprintf("service2 response: %s \n\n", <-service2Result)

		fmt.Fprintf(w, response1Message+response2Message)
	})
	http.ListenAndServe(":8080", nil)
}
