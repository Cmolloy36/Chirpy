package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerPostPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type inputJSON struct {
		Event string `json:"event,omitempty"`
		Data  struct {
			UserID uuid.UUID `json:"user_id,omitempty"`
		} `json:"data,omitempty"`
	}

	var inputData inputJSON

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&inputData); err != nil {
		errorMessage := "Something went wrong"

		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	if inputData.Event != "user.upgraded" {
		respondwithJSON(w, http.StatusNoContent, nil)
		return
	}

	dbUser, err := apiCfg.dbQueries.UpgradeUsertoChirpyRed(r.Context(), inputData.Data.UserID)
	if err != nil {
		errorMessage := err.Error()

		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, errorMessage)
			return
		}

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

	respondwithJSON(w, http.StatusNoContent, user)
}
