package main

import (
	"encoding/json"
	"net/http"

	"github.com/Cmolloy36/Chirpy/internal/auth"
	"github.com/Cmolloy36/Chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Error encountered while decoding email & password"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	dbPassword, err := apiCfg.dbQueries.GetPassword(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := "Error encountered while retrieving hashed password"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if err := auth.CheckPasswordHash(dbPassword, inputData.Password); err != nil {
		errorMessage := "Incorrect email or password"

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	dbUser, err := apiCfg.dbQueries.GetUser(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := "Error encountered when retrieving user"

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

func (apiCfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Error encountered while decoding email & password"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	hashedPassword, err := auth.HashPassword(inputData.Password)
	if err != nil {
		errorMessage := "Error encountered while hashing password"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          inputData.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := apiCfg.dbQueries.CreateUser(r.Context(), createUserParams)
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
