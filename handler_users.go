package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Cmolloy36/Chirpy/internal/auth"
	"github.com/Cmolloy36/Chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if (inputData.ExpiresInSeconds == 0) || (inputData.ExpiresInSeconds > 3600) {
		inputData.ExpiresInSeconds = 3600
	}

	dbPassword, err := apiCfg.dbQueries.GetPassword(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if err := auth.CheckPasswordHash(dbPassword, inputData.Password); err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	dbUser, err := apiCfg.dbQueries.GetUser(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := "Error encountered when retrieving user"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, apiCfg.secretString)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	currTime := time.Now()
	expiresIn := time.Duration(60*24) * time.Duration(time.Hour)
	expiresAt := currTime.Add(expiresIn)

	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    dbUser.ID,
		ExpiresAt: expiresAt,
	}

	_, err = apiCfg.dbQueries.CreateRefreshToken(r.Context(), createRefreshTokenParams)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        accessToken,
		RefreshToken: refreshTokenString,
		IsChirpyRed:  dbUser.IsChirpyRed,
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
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	hashedPassword, err := auth.HashPassword(inputData.Password)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          inputData.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := apiCfg.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondwithJSON(w, http.StatusCreated, user)

}

func (apiCfg *apiConfig) handlerPutUser(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	hashedPassword, err := auth.HashPassword(inputData.Password)
	if err != nil {
		errorMessage := err.Error()

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

	UpdateUserCredentialsParams := database.UpdateUserCredentialsParams{
		ID:             userID,
		Email:          inputData.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := apiCfg.dbQueries.UpdateUserCredentials(r.Context(), UpdateUserCredentialsParams)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondwithJSON(w, http.StatusOK, user)

}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Email string `json:"email"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	dbUser, err := apiCfg.dbQueries.GetUser(r.Context(), inputData.Email)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondwithJSON(w, http.StatusOK, user)

}

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Token struct {
		Token string `json:"token"`
	}

	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	refreshTokenParams, err := apiCfg.dbQueries.GetRefreshToken(r.Context(), refreshTokenString)
	if (err != nil) || (time.Now().After(refreshTokenParams.ExpiresAt)) {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	if !refreshTokenParams.RevokedAt.Time.IsZero() {
		errorMessage := "refresh token was previously revoked"

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	user, err := apiCfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshTokenParams.Token)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, apiCfg.secretString)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	tokenStruct := Token{
		Token: accessToken,
	}

	respondwithJSON(w, http.StatusOK, tokenStruct)
}

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorMessage := err.Error()

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	revokedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	setTokenRevokedAtParams := database.SetTokenRevokedAtParams{
		Token:     refreshTokenString,
		RevokedAt: revokedAt,
	}

	refreshTokenParams, err := apiCfg.dbQueries.SetTokenRevokedAt(r.Context(), setTokenRevokedAtParams)
	if (err != nil) || (time.Now().After(refreshTokenParams.ExpiresAt)) {
		errorMessage := err.Error()

		respondWithError(w, http.StatusUnauthorized, errorMessage)
		return
	}

	respondwithJSON(w, http.StatusNoContent, nil)
}
