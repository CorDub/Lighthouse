package handlers

import (
	"database/sql"
	"net/http"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) CheckPassword(w http.ResponseWriter, r *http.Request) {
	//extract email sent
	type body struct {
		Email string `json:"email" validate:"required,email,max=254"`
	}

	decodedBody := body{}
	err := decodeRequestBody(r, &decodedBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't decode the body", err)
		return
	}

	//validate that a user exist with this email
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), decodedBody.Email) 
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "There are no accounts with this email", err)
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	//revoke any pending magic link tokens
	errRevoke := apiCfg.DB.RevokeMagicLinkTokensFromUser(r.Context(), user.ID)
	if errRevoke != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	//send an email with new magic link token
	rawToken := auth.MakeRefreshToken()
	params := database.CreateMagicLinkTokenParams{
		Token: rawToken,
		UserID: user.ID,
	}
	magicLinkToken, err := apiCfg.DB.CreateMagicLinkToken(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	errSendingPasswordEmail := apiCfg.sendPasswordRecoveryEmail(magicLinkToken.Token, "corentindubois22@gmail.com")
	if errSendingPasswordEmail != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected error sending recovery email", errSendingPasswordEmail)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}