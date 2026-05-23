package handlers

import (
	"log"
	"net/http"
	"time"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email" validate:"required,email,max=254"`
		Password string `json:"password" validate:"required,min=8,max=72,printascii"`
	}
	type response struct {
		User
	}

	// decode params
	params := parameters{}
	err := decodeRequestBody(r, &params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// get user by email
	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// return if no password (gogle login)
	if user.HashedPassword.Valid == false {
		RespondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	
	// check password is correct
	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil || !match {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// get jwt token
	token, err := auth.MakeJWT(user.ID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// set the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
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
		Secure: apiCfg.Secure,
		SameSite: apiCfg.SameSite,
		Path: "/",
		MaxAge: 30 * 24 * 60 * 60,
	})

	handlerUser := dbUserToUser(user)

	RespondWithJSON(w, http.StatusOK, response{
		User: handlerUser,
	})
}