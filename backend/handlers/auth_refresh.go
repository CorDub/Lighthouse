package handlers

import (
	"net/http"
	"database/sql"
	"time"
	"log"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	//get refresh token
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	tokenString := cookie.Value

	// get token from db
	token, err := apiCfg.DB.GetRefreshTokenByToken(r.Context(), tokenString)

	// check if valid
	if err == sql.ErrNoRows {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	if token.RevokedAt.Valid || token.ExpiresAt.Before(time.Now()) {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//if valid, revoke the old refresh token
	_, errRevoke := apiCfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if errRevoke != nil {
		log.Printf("Couldn't revoke the refresh token: %s", errRevoke)
		return
	}

	// create a new refresh token
	newRefreshToken := auth.MakeRefreshToken()
	newRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: newRefreshToken,
		UserID: token.UserID,
	}
	
	if err := apiCfg.DB.CreateRefreshToken(r.Context(), newRefreshTokenParams); err != nil {
		log.Printf("Couldn't save the new refresh token in the DB : %s", err)
		return
	}

	// create also a new JWT cookie
	accessToken, err := auth.MakeJWT(token.UserID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the new refesh cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: newRefreshToken,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	// set the new JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: accessToken,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 15 * 60,
	})

	RespondWithJSON(w, http.StatusOK, nil)
}
