package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerCleanJSON(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Body string `json:"body"`
	}

	type outputJSON struct {
		CleanedBody string `json:"cleaned_body"`
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

	cleanedBodyJSON := outputJSON{
		CleanedBody: cleanedBody,
	}

	respondwithJSON(w, http.StatusOK, cleanedBodyJSON)

}

func (apiCfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Email string `json:"email"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Error encountered while decoding email"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	dbUser, err := apiCfg.dbQueries.CreateUser(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := "Could not create user"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondwithJSON(w, http.StatusCreated, user)

}

var profaneWords []string = []string{"kerfuffle", "sharbert", "fornax"}

func removeProfanity(chirp string) string {
	wordSlc := strings.Split(chirp, " ")
	for i, _ := range wordSlc {
		for j, _ := range profaneWords {
			if strings.ToLower(wordSlc[i]) == profaneWords[j] {
				wordSlc[i] = "****"
			}
		}
	}

	cleanSlc := strings.Join(wordSlc, " ")

	return cleanSlc
}
