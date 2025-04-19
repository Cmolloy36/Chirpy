package main

import (
	"net/http"
)

func main() {
	newServeMux := &http.ServeMux{}

	newHttpServer := http.Server{
		Addr:    ":8080",
		Handler: newServeMux,
	}

	newServeMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	newServeMux.HandleFunc("/", handler)

	newHttpServer.ListenAndServe()

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
