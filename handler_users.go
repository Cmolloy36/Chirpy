package main

import (
	"encoding/json"
	"net/http"
)

func (apiCfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
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

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
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

	dbUser, err := apiCfg.dbQueries.GetUser(r.Context(), inputData.Email)
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

	respondwithJSON(w, http.StatusOK, user)

}
