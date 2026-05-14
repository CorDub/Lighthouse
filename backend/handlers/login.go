package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// dev settings
	sameSite := http.SameSiteStrictMode
	secure := true

	if apiCfg.Env == "dev" {
		sameSite = http.SameSiteLaxMode
		secure = false
	}

	// set the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		Secure: secure,
		SameSite: sameSite,
		Path: "/",
		MaxAge: 15 * 60,
	})

	// create refresh token
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}

	// save it in DB
	if err := apiCfg.DB.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		log.Printf("Couldn't save the refresh token in db: %s", err)
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the refresh cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		HttpOnly: true,
		Secure: secure,
		SameSite: sameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}


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

	// set both cookies	
	sameSite := http.SameSiteStrictMode
	secure := true

	if apiCfg.Env == "dev" {
		sameSite = http.SameSiteLaxMode
		secure = false
	}

	// set the new refesh cookie
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: newRefreshToken,
		HttpOnly: true,
		Secure: secure,
		SameSite: sameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	// set the new JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: accessToken,
		HttpOnly: true,
		Secure: secure,
		SameSite: sameSite,
		Path: "/",
		MaxAge: 15 * 60,
	})

	RespondWithJSON(w, http.StatusOK, nil)
}


func (apiCfg *ApiConfig) Logout(w http.ResponseWriter, r *http.Request) {
	//get the refreshToken value
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	tokenString := cookie.Value

	//set refresh and JWT tokens to MaxAge -1 to clear them in the browser
	http.SetCookie(w, &http.Cookie{
		Name: "refreshToken",
		Value: "",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
		Path: "/",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: "",
		HttpOnly: true,
		Expires: time.Unix(0, 0),
		Path: "/",
		MaxAge: -1,
	})

	// revoke the refresh token in the DB
	_, errRevoke := apiCfg.DB.RevokeRefreshToken(r.Context(), tokenString)
	if errRevoke != nil {
		log.Printf("Couldn't revoke the  refresh token: %s", errRevoke)
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", errRevoke)
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}