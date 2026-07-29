package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"fmt"

	"github.com/google/uuid"

	"Lighthouse/internal/database"
	"Lighthouse/internal/auth"
)

func (apiCfg *ApiConfig) CreateMagicLinkInvite (w http.ResponseWriter, r *http.Request) {
	// decode body for user id
	type Body struct {
		ID uuid.UUID `json:"userId" validate:"required"`
		Name string `json:"name" validate:"required"`
	}
	decodedBody := Body{}

	err := decodeRequestBody(r, &decodedBody)
	if (err != nil) {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode the body of the request", err)
		return
	}

	// validate
	errValidation := validate.Struct(decodedBody)
	if errValidation != nil {
		RespondWithError(w, http.StatusInternalServerError, "Invalid user Id provided", errValidation)
		return
	}

	// create ML Token
	token := auth.MakeRefreshToken()
	params := database.CreateMagicLinkTokenParams{
		Token: token,
		UserID: decodedBody.ID,
		Name: sql.NullString{String: decodedBody.Name, Valid: true},
	}

	mlToken, err := apiCfg.DB.CreateMagicLinkToken(r.Context(), params)
	if (err != nil) {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create the magic link token", err)
		return
	}

	link := fmt.Sprintf("http://localhost:5173/invite?token=%s&name=%s", mlToken.Token, url.QueryEscape(decodedBody.Name))

	RespondWithJSON(w, http.StatusCreated, link)
}