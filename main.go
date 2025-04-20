package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	newServeMux := &http.ServeMux{}

	newHttpServer := http.Server{
		Addr:    ":8080",
		Handler: newServeMux,
	}

	apiCfg := apiConfig{}

	funcHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	newServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(funcHandler))

	newServeMux.HandleFunc("GET /api/healthz", handler)

	newServeMux.HandleFunc("POST /api/validate_chirp", handlerJSON)

	newServeMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	newServeMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	newHttpServer.ListenAndServe()

}

type apiConfig struct {
	fileserverHits atomic.Int32 // allows us to safely increment & read across multiple goroutines
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handlerJSON(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Body string `json:"body"`
	}

	type errStruct struct {
		ErrorMessage string `json:"error"`
	}

	type validStruct struct {
		Validity bool `json:"valid"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorStruct := errStruct{
			ErrorMessage: "Something went wrong",
		}

		dat, err := json.Marshal(errorStruct)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Write(dat)
		return
	}

	if len(inputData.Body) > 140 {
		w.WriteHeader(400)
		errorStruct := errStruct{
			ErrorMessage: "Chirp is too long",
		}

		dat, err := json.Marshal(errorStruct)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Write(dat)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	ret := validStruct{
		Validity: true,
	}

	dat, err := json.Marshal(ret)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Write(dat)
}

func (apiCfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html`, apiCfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (apiCfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Write([]byte("Hits: 0"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
