package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Cmolloy36/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %w", err))
	}

	dbQueries := database.New(db)

	newServeMux := &http.ServeMux{}

	newHttpServer := http.Server{
		Addr:    ":8080",
		Handler: newServeMux,
	}

	apiCfg := apiConfig{}
	apiCfg.platform = os.Getenv("PLATFORM")
	apiCfg.fileserverHits.Store(0)
	apiCfg.dbQueries = dbQueries
	apiCfg.secretString = os.Getenv("SIGNING_SECRET")
	apiCfg.polkaKey = os.Getenv("POLKA_KEY")

	funcHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	newServeMux.Handle("/app/", apiCfg.middlewareMetricsInc(funcHandler))

	newServeMux.HandleFunc("GET /api/healthz", handler)

	newServeMux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)

	newServeMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	newServeMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	newServeMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)

	newServeMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	newServeMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPostPolkaWebhook)

	newServeMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	newServeMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	newServeMux.HandleFunc("POST /api/users", apiCfg.handlerPostUser)

	newServeMux.HandleFunc("PUT /api/users", apiCfg.handlerPutUser)

	newServeMux.HandleFunc("GET /api/users/{userID}", apiCfg.handlerGetUser)

	newServeMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	newServeMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	newHttpServer.ListenAndServe()

}

type apiConfig struct {
	platform       string
	fileserverHits atomic.Int32 // allows us to safely increment & read across multiple goroutines
	dbQueries      *database.Queries
	secretString   string
	polkaKey       string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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
	// apiCfg.fileserverHits.Store(0)
	// w.Write([]byte("Hits: 0"))

	if apiCfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	apiCfg.dbQueries.ResetUsers(context.Background())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
