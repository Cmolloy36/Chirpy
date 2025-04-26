package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Cmolloy36/Chirpy/internal/auth"
	"github.com/Cmolloy36/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		errorMessage := "Error parsing chirp ID"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	accessTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	userID, err := auth.ValidateJWT(accessTokenString, apiCfg.secretString)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		errorMessage := err.Error()

		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, errorMessage)
			return
		}

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if chirp.UserID != userID {
		errorMessage := "not authorized to delete this Chirp"

		respondWithError(w, http.StatusForbidden, errorMessage)
		return
	}

	if err := apiCfg.dbQueries.DeleteChirp(r.Context(), chirpID); err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	respondwithJSON(w, http.StatusNoContent, nil)
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
		errorMessage := err.Error()

		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, errorMessage)
			return
		}

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
	s := r.URL.Query().Get("author_id")
	// s is a string that contains the value of the author_id query parameter
	// if it exists, or an empty string if it doesn't

	userID, err := uuid.Parse(s)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	var chirpSlc []database.Chirp

	if s != "" {
		chirpSlc, err = apiCfg.dbQueries.GetChirpsForUser(context.Background(), userID)
	} else {
		chirpSlc, err = apiCfg.dbQueries.GetChirps(context.Background())
	}

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

func (apiCfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Body string `json:"body"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Something went wrong"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	validatedUserID, err := auth.ValidateJWT(token, apiCfg.secretString)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
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
		UserID: validatedUserID,
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
