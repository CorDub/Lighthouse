package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"fmt"
	"context"
	"errors"

	"github.com/google/uuid"

	"Lighthouse/internal/database"
	"Lighthouse/internal/auth"
)

type createMagicLinkInviteParams struct {
	ID uuid.UUID `json:"userId" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (apiCfg *ApiConfig) CreateMagicLinkInvite (w http.ResponseWriter, r *http.Request) {
	// decode body for user id
	decodedBody := createMagicLinkInviteParams{}

	err := decodeRequestBody(r, &decodedBody)
	if (err != nil) {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode the body of the request", err)
		return
	}

	createdMagicLinkToken, err := apiCfg.createMagicLinkInvite(r.Context(), decodedBody)
	errCreatingMagicLinkToken := apiCfg.respondCreateMagicLinkInviteError(w, err)
	if errCreatingMagicLinkToken != nil {
		return
	}

	link := fmt.Sprintf("http://localhost:5173/invite?token=%s&name=%s", createdMagicLinkToken.Token, url.QueryEscape(decodedBody.Name))

	RespondWithJSON(w, http.StatusCreated, link)
}


func (apicfg *ApiConfig) respondCreateMagicLinkInviteError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrNotValidated):
		RespondWithError(w, http.StatusInternalServerError, "Invalid user Id provided", err)
	case errors.Is(err, ErrCreatingMagicLinkToken):
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create the magic link token", err)
	default:
		RespondWithError(w, http.StatusInternalServerError, "Unexpected error", err) 
	}
	return err
}


func (apiCfg *ApiConfig) createMagicLinkInvite (
	ctx context.Context, decodedBody createMagicLinkInviteParams,
	) (database.MagicLinkToken, error) {
	// validate
	errValidation := validate.Struct(decodedBody)
	if errValidation != nil {
		return database.MagicLinkToken{}, ErrNotValidated
	}

	// create ML Token
	token := auth.MakeRefreshToken()
	params := database.CreateMagicLinkTokenParams{
		Token: token,
		UserID: decodedBody.ID,
		Name: sql.NullString{String: decodedBody.Name, Valid: true},
	}

	mlToken, err := apiCfg.DB.CreateMagicLinkToken(ctx, params)
	if (err != nil) {
		return database.MagicLinkToken{}, ErrCreatingMagicLinkToken
	}

	return mlToken, nil
}