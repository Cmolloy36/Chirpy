package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Cmolloy36/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Something went wrong"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if len(inputData.Body) > 140 {

		errorMessage := "Chirp is too long"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	cleanedBody := removeProfanity(inputData.Body)

	createChirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: inputData.UserID,
	}

	chirp, err := apiCfg.dbQueries.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		errorMessage := "Error creating chirp"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	retChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      cleanedBody,
		UserID:    chirp.UserID,
	}

	respondwithJSON(w, http.StatusCreated, retChirp)
}

func (apiCfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		errorMessage := "Error parsing chirp ID"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirp(context.Background(), chirpID)
	if err != nil {
		errorMessage := "Error getting chirp"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	retChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondwithJSON(w, http.StatusOK, retChirp)
}

func (apiCfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpSlc, err := apiCfg.dbQueries.GetChirps(context.Background())
	if err != nil {
		errorMessage := "Error getting chirps"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	retSlc := make([]Chirp, len(chirpSlc))

	for i, chirp := range chirpSlc {
		retSlc[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondwithJSON(w, http.StatusOK, retSlc)
}

var profaneWords []string = []string{"kerfuffle", "sharbert", "fornax"}

func removeProfanity(chirpBody string) string {
	wordSlc := strings.Split(chirpBody, " ")
	for i, _ := range wordSlc {
		for j, _ := range profaneWords {
			if strings.ToLower(wordSlc[i]) == profaneWords[j] {
				wordSlc[i] = "****"
			}
		}
	}

	cleanStr := strings.Join(wordSlc, " ")

	return cleanStr
}
