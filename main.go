package main

import (
	"net/http"
)

func main() {
	newServeMux := http.ServeMux{}

	newHttpServer := http.Server{
		Addr:    ":8080",
		Handler: &newServeMux,
	}

	newHttpServer.ListenAndServe()
}
