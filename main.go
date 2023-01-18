package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request url:", r.RequestURI, "request path:", r.URL.Path)
	})

	log.Println("Listening :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln("something wrong:", err)
	}
}
