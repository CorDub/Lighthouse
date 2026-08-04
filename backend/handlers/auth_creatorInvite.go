package handlers

import (
	"net/http"
	"time"
	"log"

	"Lighthouse/internal/auth"
	"Lighthouse/internal/database"
)

func (apiCfg *ApiConfig) HandleUserInvite(w http.ResponseWriter, r *http.Request) {
	//createUserParams defined in users_create
	decodedBody := createUserParams{}

	err:= decodeRequestBody(r, &decodedBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode the body of the request", err)
		return
	}

	//create user - validation inside this function
	createdUser, err := apiCfg.createUser(r.Context(), decodedBody)
	errCheckError := apiCfg.respondCreateUserError(w, err) 
	if errCheckError != nil {
		return
	}

	//get jwt token
	token, err := auth.MakeJWT(createdUser.ID, apiCfg.JWT, time.Minute * 15)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unexpected server error", err)
		return
	}

	// create refresh token
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: createdUser.ID,
	}

	// save it in DB
	if err := apiCfg.DB.CreateRefreshToken(r.Context(), refreshTokenParams); err != nil {
		log.Printf("Couldn't save the refresh token in db: %s", err)
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

	// respond
	w.WriteHeader(http.StatusCreated)
}