package handlers

import (
	"net/http"

	"Lighthouse/internal/database"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) UpdateUserLanguage (w http.ResponseWriter, r *http.Request) {
	//validate, decode body and get params
	type Body struct {
		Language string `json:"language" validate:"required,oneof=en es"`
	}
	decodedBody := Body{}
	err := decodeRequestBody(r, &decodedBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode the body of the request", err)
		return
	}

	strUserId := r.PathValue("id")
	uuidUserId, err := uuid.Parse(strUserId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not parse the userId", err)
		return
	}

	// patch the user in DB
	params := database.UpdateUserLanguageParams{
		ID: uuidUserId,
		Language: decodedBody.Language,
	}

	errUpdate := apiCfg.DB.UpdateUserLanguage(r.Context(), params)
	if errUpdate != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not updqte the user", err)
		return
	}

	// confirm
	RespondWithJSON(w, http.StatusNoContent, nil)
}
