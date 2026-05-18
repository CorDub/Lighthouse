package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"Lighthouse/internal/auth"
)

func (apiCfg *ApiConfig) ChangePassword(w http.ResponseWriter, r *http.Request) {
	//decode body
	type body struct {
		Password string `json:"password" validate:"required,min=8,max=72,printascii"`
		Token string `json:"token" validate:"required,len=64"`
	}
	decodedBody := body{}
	err := decodeRequestBody(r, &decodedBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	//check that the token is valid
	dbToken, err := apiCfg.DB.GetMagicLinkTokenByToken(r.Context(), decodedBody.Token)
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if time.Now().UTC().After(dbToken.ExpiresAt) {
		fmt.Printf("expires_at: %v", dbToken.ExpiresAt)
		fmt.Printf("time.Now.UTC :%v", time.Now().UTC())
		RespondWithError(w, http.StatusBadRequest, "Reset link has expired", fmt.Errorf("Reset link has expired"))
		return
	}

	//check that the password is not the previous password
	user, err := apiCfg.DB.GetUserByID(r.Context(), dbToken.UserID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}
	
	if user.HashedPassword.Valid {
		match, err := auth.CheckPasswordHash(decodedBody.Password, user.HashedPassword.String)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
			return
		}

		if match {
			RespondWithError(w, http.StatusBadRequest, "Same password", errors.New("Same password"))
			return
		}
	}

	//revoke the token
	errRevoke := apiCfg.DB.RevokeMagicLinkToken(r.Context(), dbToken.Token)
	if errRevoke != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	//send confirmation
	RespondWithJSON(w, http.StatusNoContent, nil)
}
