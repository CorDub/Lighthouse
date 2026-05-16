package handlers

import (
	"net/http"

	"Lighthouse/internal/auth"
)

func (apiCfg *ApiConfig) CheckAuth(w http.ResponseWriter, r *http.Request) {
	type respAuth struct {
		User User `json:"user"`
	}

	//get JWT token
	cookie, err := r.Cookie("token")

	// if you don't have the JWT check if the refresh token is present
	if err == http.ErrNoCookie {
		apiCfg.Refresh(w, r)
		return
	}
	
	// if it's there check the validity
	tokenString := cookie.Value
	userId, err := auth.ValidateJWT(tokenString, apiCfg.JWT)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//get the user
	dbUser, err := apiCfg.DB.GetUserByID(r.Context(), userId)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	//prepare user to correct format for frontend
	handlersUser := dbUserToUser(dbUser)

	// prepare and send response
	userPayload := respAuth{
		User: handlersUser,
	}
	RespondWithJSON(w, http.StatusAccepted, userPayload)
}